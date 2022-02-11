package users_api

import (
	"git.thorn.sh/Thorn/go-vodstream/cmd/users_api/handlers"
	"git.thorn.sh/Thorn/go-vodstream/storage/sql"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
)

// NewCmd registers the cobra command to be called from the CLI.
func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "users_api",
		Short: "launches an auth/users API provider.",
		Run:   Start,
	}
}

// Start the http(s) listen server.
func Start(_ *cobra.Command, _ []string) {
	log.Info("Starting Users API")

	bindAddress := viper.GetString("api.users.bind_address")

	db := sql.NewDbConn()
	storage := sql.NewUserStorage(db)

	handler := handlers.NewAuthHandler(storage)

	r := mux.NewRouter()
	handler.RegisterRoutes(r)

	log.Println("Listening on" + bindAddress + "..")
	err := http.ListenAndServe(bindAddress, r)
	if err != nil {
		log.Fatal(err)
	}
}
