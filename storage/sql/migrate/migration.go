package migrate

import (
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/fs"
	"log"
	"regexp"
)

// MigrationFS contains all the SQL migration files.
//go:embed migrations/*.sql
var MigrationFS embed.FS

const migrationsTableName = "migrations"

type Migration struct {
	TableName     string
	MigrationName string
	UpQuery       string
	DownQuery     string
	db            *sqlx.DB
}

type MigrationType int

const (
	UpMigration = iota
	DownMigration
)

func (m Migration) Up() error {
	// STUB
	return nil
}

func (m Migration) Down() error {
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

func getMigrationsForTable(dir []fs.DirEntry, tableName string) []Migration {
	migrations := make([]Migration, 0)

	/*for _, dirEntry := range dir {
		if dirEntry.IsDir() {
			continue
		}

		valid, migrationType := validMigrationFile(dirEntry.Name(), tableName)
		if !valid {
			continue
		}


	} */

	// STUB
	return migrations
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
	sql := `SELECT MAX(batch) from $1`

	row := db.QueryRow(sql, migrationsTableName)

	var batch int
	err := row.Scan(&batch)
	if err != nil {
		log.Fatal(err)
	}

	return batch
}

// shouldRunMigration returns true if a migration by this name hasn't been run
// in a previous batch.
func shouldRunMigration(db *sqlx.DB, migrationName string) bool {
	sql := `SELECT id FROM $1 WHERE migration = $2`

	row := db.QueryRow(sql, migrationsTableName, migrationName)

	var id int
	err := row.Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	return id >= 0
}

// migrationTableExists returns true if the migrations table exists in the given database.
func migrationTableExists(db *sqlx.DB) bool {
	sql := `SELECT EXISTS (
		SELECT FROM 
			pg_tables
		WHERE 
			schemaname = 'public' AND 
			tablename  = '$1'
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
