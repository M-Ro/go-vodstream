package migrate

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"os"
	"testing"
	"time"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)
	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sqlx.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Run tests
	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestMigration_validMigrationFile(t *testing.T) {
	tests := []struct {
		testName  string
		fileName  string
		tableName string

		expectedValid bool
		expectedType  MigrationType
	}{
		{
			testName:      "expect fail with empty fileName",
			fileName:      "",
			tableName:     "users",
			expectedValid: false,
			expectedType:  DownMigration,
		},
		{
			testName:      "expect fail with empty tableName and empty fileName",
			fileName:      "",
			tableName:     "",
			expectedValid: false,
			expectedType:  DownMigration,
		},
		{
			testName:      "expect fail with incorrect filename format",
			fileName:      "users202008101432_left.sql",
			tableName:     "users",
			expectedValid: false,
			expectedType:  DownMigration,
		},
		{
			testName:      "expect fail with migration file for wrong table",
			fileName:      "users_202202041432_up.sql",
			tableName:     "broadcasts",
			expectedValid: false,
			expectedType:  DownMigration,
		},
		{
			testName:      "expect pass, UpMigration",
			fileName:      "table_202202041432_create_table_up.sql",
			tableName:     "table",
			expectedValid: true,
			expectedType:  UpMigration,
		},
		{
			testName:      "expect pass, DownMigration",
			fileName:      "table_202202041432_create_table_down.sql",
			tableName:     "table",
			expectedValid: true,
			expectedType:  DownMigration,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			valid, migrationType := validMigrationFile(test.fileName, test.tableName)

			// Compare output to expected output.
			if !cmp.Equal(valid, test.expectedValid) {
				t.Fatal(cmp.Diff(valid, test.expectedValid))
			}

			if !cmp.Equal(migrationType, test.expectedType) {
				t.Fatal(cmp.Diff(migrationType, test.expectedType))
			}
		})
	}
}

func TestMigration_getMigrationsForTable(t *testing.T) {
	type mockfile struct {
		name  string
		isDir bool
	}

	tests := []struct {
		testName    string
		tableName   string
		dirContents []mockfile

		expectedReturn []Migration
		expectedErr    error
	}{
		{
			testName:    "expect nothing with empty directory",
			tableName:   "users",
			dirContents: []mockfile{},

			expectedReturn: []Migration{},
			expectedErr:    nil,
		},
		{
			testName:  "expect nothing with no migrations for specified table",
			tableName: "users",
			dirContents: []mockfile{
				{
					name:  "tests_202205132200_create_table_up.sql",
					isDir: false,
				},
				{
					name:  "tests_202205132200_create_table_down.sql",
					isDir: false,
				},
			},

			expectedReturn: []Migration{},
		},
		{
			testName:  "expect nothing with various invalid files/dirs",
			tableName: "users",
			dirContents: []mockfile{
				{
					name:  "tests_202205132200_create_table_up.sql",
					isDir: false,
				},
				{
					name:  "users_202205132200_create_table_up.sql",
					isDir: true,
				},
				{
					name:  "users_202205132200_create_table_down.sql",
					isDir: true,
				},
				{
					name:  "users_202105132200_down.sql",
					isDir: false,
				},
				{
					name:  "users_create_table_down.sql",
					isDir: false,
				},
			},

			expectedReturn: []Migration{},
			expectedErr:    nil,
		},
		{
			testName:  "expect error with up/down migration missing counter component",
			tableName: "users",
			dirContents: []mockfile{
				{
					name:  "users_202205132200_create_table_down.sql",
					isDir: false,
				},
			},

			expectedReturn: []Migration{},
			expectedErr:    ErrHalfMigration,
		},
		{
			testName:  "expect array(2) of valid migrations",
			tableName: "testtable",
			dirContents: []mockfile{
				{
					name:  "testtable_202205132200_create_table_down.sql",
					isDir: false,
				},
				{
					name:  "testtable_202205132200_create_table_up.sql",
					isDir: false,
				},
				{
					name:  "testtable_202205132300_update_table_down.sql",
					isDir: false,
				},
				{
					name:  "testtable_202205132300_update_table_up.sql",
					isDir: false,
				},
			},

			expectedReturn: []Migration{
				{
					TableName:     "testtable",
					MigrationName: "testtable_202205132200_create_table",
					UpQuery:       "sql here",
					DownQuery:     "sql here",
				},
				{
					TableName:     "testtable",
					MigrationName: "testtable_202205132300_update_table",
					UpQuery:       "sql here",
					DownQuery:     "sql here",
				},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create mock filesystem
			var fs = afero.NewMemMapFs()
			for _, f := range test.dirContents {
				if f.isDir {
					fs.Mkdir(f.name, 0755)
				} else {
					afero.WriteFile(fs, f.name, []byte("sql here"), 0755)
				}
			}

			files, err := afero.ReadDir(fs, ".")
			if err != nil {
				t.Fatal(err)
			}

			migrations, err := getMigrationsForTable(fs, files, test.tableName)

			if !cmp.Equal(err, test.expectedErr, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedErr, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(migrations, test.expectedReturn) {
				t.Fatal(cmp.Diff(migrations, test.expectedReturn))
			}
		})
	}
}

func TestMigration_createMigrationTable(t *testing.T) {
	createMigrationTable(db)

	row := db.QueryRow(`SELECT EXISTS (
		SELECT FROM 
			pg_tables
		WHERE 
			schemaname = 'public' AND 
			tablename  = $1
		);`, migrationsTableName)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("table not found")
	}

	// Delete the table again for next tests
	delQuery := fmt.Sprintf(`DROP TABLE %s`, migrationsTableName)
	db.MustExec(delQuery)
}

func TestMigration_migrationTableExists(t *testing.T) {
	// Should not exist
	exists := migrationTableExists(db)
	if exists {
		t.Fatal("table found when not expected")
	}

	// Create and run test again expecting to return true
	createMigrationTable(db)

	exists = migrationTableExists(db)
	if !exists {
		t.Fatal("table not found when expected")
	}

	// Delete the table again for next tests
	delQuery := fmt.Sprintf(`DROP TABLE %s`, migrationsTableName)
	db.MustExec(delQuery)
}

func TestMigration_lastBatchNumber(t *testing.T) {
	type migrationRow struct {
		migration string
		batch     int
	}

	tests := []struct {
		testName      string
		tableContents []migrationRow
		expectedValue int
	}{
		{
			testName:      "expect 0 with empty table",
			tableContents: []migrationRow{},
			expectedValue: 0,
		},
		{
			testName: "expect 1",
			tableContents: []migrationRow{
				{
					migration: "tests_202202062251_create_table",
					batch:     1,
				},
			},
			expectedValue: 1,
		},
		{
			testName: "expect 4 even with gaps",
			tableContents: []migrationRow{
				{
					migration: "tests_202202062251_create_table",
					batch:     1,
				},
				{
					migration: "users_202202062251_update_table",
					batch:     4,
				},
			},
			expectedValue: 4,
		},
	}

	createMigrationTable(db)

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Ensure table is empty
			db.MustExec(fmt.Sprintf("DELETE FROM %s", migrationsTableName))

			// Add our rows
			for _, row := range test.tableContents {
				insQuery := fmt.Sprintf(`INSERT INTO %s (migration, batch) VALUES ($1, $2)`, migrationsTableName)
				db.MustExec(insQuery, row.migration, row.batch)
			}

			// Call lastBatchNumber and check number
			value := lastBatchNumber(db)

			if !cmp.Equal(value, test.expectedValue) {
				t.Fatal(cmp.Diff(value, test.expectedValue))
			}
		})
	}

	// Delete the table again for next tests
	delQuery := fmt.Sprintf(`DROP TABLE %s`, migrationsTableName)
	db.MustExec(delQuery)
}

func TestMigration_shouldRunMigration(t *testing.T) {
	type migrationRow struct {
		migration string
		batch     int
	}

	tests := []struct {
		testName      string
		tableContents []migrationRow
		migrationName string
		expectedValue bool
	}{
		{
			testName:      "expect true with empty table",
			tableContents: []migrationRow{},
			migrationName: "table_202202021453_create_table",
			expectedValue: true,
		},
		{
			testName: "expect true with migration not in table",
			tableContents: []migrationRow{
				{
					migration: "tests_202202062251_create_table",
					batch:     1,
				},
			},
			migrationName: "tests_202203062251_update_table",
			expectedValue: true,
		},
		{
			testName: "expect false with existing migration",
			tableContents: []migrationRow{
				{
					migration: "tests_202202062251_create_table",
					batch:     1,
				},
			},
			migrationName: "tests_202202062251_create_table",
			expectedValue: false,
		},
	}

	createMigrationTable(db)

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Ensure table is empty
			db.MustExec(fmt.Sprintf("DELETE FROM %s", migrationsTableName))

			// Add our rows
			for _, row := range test.tableContents {
				insQuery := fmt.Sprintf(`INSERT INTO %s (migration, batch) VALUES ($1, $2)`, migrationsTableName)
				db.MustExec(insQuery, row.migration, row.batch)
			}

			// Call lastBatchNumber and check number
			value := shouldRunMigration(db, test.migrationName)

			if !cmp.Equal(value, test.expectedValue) {
				t.Fatal(cmp.Diff(value, test.expectedValue))
			}
		})
	}

	// Delete the table again for next tests
	delQuery := fmt.Sprintf(`DROP TABLE %s`, migrationsTableName)
	db.MustExec(delQuery)
}
