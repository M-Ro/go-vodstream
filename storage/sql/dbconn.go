package sql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var db *sqlx.DB

func getConfig() {

}

func getPostgresDSN() string {
	host := viper.GetString("postgres.host")
	port := viper.GetString("postgres.port")
	user := viper.GetString("postgres.user")
	pass := viper.GetString("postgres.password")
	dbname := viper.GetString("postgres.database")

	return fmt.Sprintf(
		`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
		host, port, user, pass, dbname,
	)

	return fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		user, pass, host, port, dbname,
	)
}

func NewDbConn() *sqlx.DB {
	if db != nil {
		return db
	}

	newDb, err := sqlx.Open("postgres", getPostgresDSN())
	if err != nil {
		log.Fatal(err)
	}

	db = newDb

	return db
}
