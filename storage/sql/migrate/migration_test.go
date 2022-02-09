package migrate

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

var db *sqlx.DB

type Mockfile struct {
	Name  string
	isDir bool
}

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
	pool.MaxWait = 60 * time.Second
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
			testName:      "Expect fail with empty fileName",
			fileName:      "",
			tableName:     "users",
			expectedValid: false,
			expectedType:  DownMigration,
		},
		{
			testName:      "Expect fail with empty tableName and empty fileName",
			fileName:      "",
			tableName:     "",
			expectedValid: false,
			expectedType:  DownMigration,
		},
		{
			testName:      "Expect fail with incorrect filename format",
			fileName:      "users202008101432_left.sql",
			tableName:     "users",
			expectedValid: false,
			expectedType:  DownMigration,
		},
		{
			testName:      "Expect fail with migration file for wrong table",
			fileName:      "users_202202041432_up.sql",
			tableName:     "broadcasts",
			expectedValid: false,
			expectedType:  DownMigration,
		},
		{
			testName:      "Expect pass, UpMigration",
			fileName:      "table_202202041432_create_table_up.sql",
			tableName:     "table",
			expectedValid: true,
			expectedType:  UpMigration,
		},
		{
			testName:      "Expect pass, DownMigration",
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
	t.Log("test not implemented")
	t.Fail()
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
	t.Log("test not implemented")
	t.Fail()
}

func TestMigration_lastBatchNumber(t *testing.T) {
	t.Log("test not implemented")
	t.Fail()
}

func TestMigration_shouldRunMigration(t *testing.T) {
	t.Log("test not implemented")
	t.Fail()
}
