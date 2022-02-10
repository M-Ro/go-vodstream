package migrate

import (
	sql2 "database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"
	"log"
	"os"
	"regexp"
)

// MigrationFS contains all the SQL migration files.
//go:embed migrations/*.sql
var MigrationFS embed.FS

var (
	ErrHalfMigration = errors.New("migration is missing up/down component")
)

const migrationsTableName = "migrations"

type Migration struct {
	TableName     string
	MigrationName string
	UpQuery       string
	DownQuery     string
}

type MigrationType int

const (
	UpMigration = iota
	DownMigration
)

func (m Migration) Up(tx *sqlx.Tx) error {
	// STUB
	return nil
}

func (m Migration) Down(tx *sqlx.Tx) error {
	// STUB
	return nil
}

func createMigrationTable(db *sqlx.DB) {
	sql := fmt.Sprintf(`CREATE TABLE %s (
				id          SERIAL   PRIMARY KEY,
				migration   TEXT     NOT NULL,
				batch       INTEGER  NOT NULL
		)`, migrationsTableName)

	db.MustExec(sql)
}

func getMigrationsForTable(fs afero.Fs, files []os.FileInfo, tableName string) ([]Migration, error) {
	strPattern := fmt.Sprintf(`^(%s_\d{12}_.+)_(down|up)\.sql$`, tableName)
	pattern := regexp.MustCompile(strPattern)

	migrations := make([]Migration, 0)

	for _, fileinfo := range files {
		if fileinfo.IsDir() {
			continue
		}

		valid, migrationType := validMigrationFile(fileinfo.Name(), tableName)
		if !valid {
			continue
		}

		match := pattern.FindStringSubmatch(fileinfo.Name())
		if match == nil || len(match) != 3 {
			continue
		}

		migrationName := match[1]
		found := false

		// Read file contents
		content, err := afero.ReadFile(fs, fileinfo.Name())
		if err != nil {
			log.Fatal(err)
		}

		// search for an existing migration and fill the missing up/down query
		for index, v := range migrations {
			if v.MigrationName == migrationName {
				found = true

				if migrationType == UpMigration {
					migrations[index].UpQuery = string(content)
				} else {
					migrations[index].DownQuery = string(content)
				}
			}
		}

		// if no migration found, add a new one to the list
		if !found {
			newMigration := Migration{
				TableName:     tableName,
				MigrationName: migrationName,
			}

			if migrationType == UpMigration {
				newMigration.UpQuery = string(content)
			} else {
				newMigration.DownQuery = string(content)
			}

			migrations = append(migrations, newMigration)
		}
	}

	// verify all migrations have an up and down query.
	for _, v := range migrations {
		if v.UpQuery == "" || v.DownQuery == "" {
			return []Migration{}, ErrHalfMigration
		}
	}

	return migrations, nil
}

// runMigration runs the given migration. The SQL called (Up or Down) is decided
// by the migrationType parameter.
func runMigration(migration Migration, migrationType MigrationType) {

}

// validMigrationFile checks if the given filename is a valid migration file for the
// given database table. Returns true and a MigrationType (up/down) on success.
// returns false on fail, and the MigrationType return should be ignored.
func validMigrationFile(filename string, tableName string) (bool, MigrationType) {
	if len(tableName) == 0 {
		return false, DownMigration
	}

	strPattern := fmt.Sprintf(`^%s_\d{12}_.+_(down|up)\.sql$`, tableName)
	pattern := regexp.MustCompile(strPattern)

	fileMatches := pattern.MatchString(filename)
	if !fileMatches {
		return false, DownMigration
	}

	upExpected := "_up.sql"
	if string(filename[len(filename)-len(upExpected):]) == upExpected {
		return true, UpMigration
	}

	return true, DownMigration
}

// lastBatchNumber returns the highest value in the batch column of the migrations table
func lastBatchNumber(db *sqlx.DB) int {
	sql := fmt.Sprintf(`SELECT MAX(batch) from %s`, migrationsTableName)

	row := db.QueryRow(sql)

	var batch sql2.NullInt32
	err := row.Scan(&batch)
	if err != nil {
		log.Fatal(err)
	}

	if !batch.Valid {
		return 0
	}

	return (int)(batch.Int32)
}

// shouldRunMigration returns true if a migration by this name hasn't been run
// in a previous batch.
func shouldRunMigration(db *sqlx.DB, migrationName string) bool {
	sql := fmt.Sprintf(`SELECT id FROM %s WHERE migration = $1 LIMIT 1`, migrationsTableName)

	row := db.QueryRow(sql, migrationName)

	var id sql2.NullInt32
	err := row.Scan(&id)

	if err == sql2.ErrNoRows {
		return true
	}

	// Any other error is fatal
	if err != nil {
		log.Fatal(err)
	}

	return !id.Valid
}

// migrationTableExists returns true if the migrations table exists in the given database.
func migrationTableExists(db *sqlx.DB) bool {
	sql := `SELECT EXISTS (
		SELECT FROM 
			pg_tables
		WHERE 
			schemaname = 'public' AND 
			tablename  = $1
		)`

	row := db.QueryRow(sql, migrationsTableName)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	return exists
}

func Run(db *sqlx.DB, tableName string) {
	// STUB
}
