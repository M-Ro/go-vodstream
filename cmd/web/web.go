package web

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
)

// NewCmd registers the cobra command to be called from the CLI.
func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "web",
		Short: "launches a web provider to deliver static FE assets",
		Run:   Start,
	}
}

// Start the http(s) listen server.
func Start(_ *cobra.Command, _ []string) {
	log.Info("Starting web frontend")

	if !FilesExist() {
		log.Fatal("Missing files for web delivery. " +
			"Ensure static directory contains compiled js/css assets from frontend repository.")
	}

	bindAddress := viper.GetString("web.bind_address")

	r := mux.NewRouter()

	r.Handle("/{file}", FsHandler())
	r.HandleFunc("/", IndexHandler)

	log.Println("Listening on" + bindAddress + "..")
	err := http.ListenAndServe(bindAddress, r)
	if err != nil {
		log.Fatal(err)
	}
}
