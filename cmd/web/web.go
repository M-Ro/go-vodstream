package web

import (
	"embed"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewCmd registers the cobra command.
func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "web",
		Short: "launches a web provider to deliver static FE assets",
		Run:   Start,
	}
}

//go:embed static
var content embed.FS

func Start(_ *cobra.Command, _ []string) {
	log.Info("Starting web frontend")

	if !filesExist(content) {
		log.Fatal("Missing files for web delivery. " +
			"Ensure static directory contains compiled js/css artifacts from frontend repository.")
	}

}

// Remove element from string slice, order not preserved
func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// FilesExist returns false if the minimum files required for web frontend delivery
// have not been compiled into the web binary.
func filesExist(fs embed.FS) bool {
	missingFiles := []string{"index.html", "app.js", "app.css"}

	files, err := fs.ReadDir("static")
	if err != nil {
		log.Error(err.Error())
		return false
	}

	for _, filename := range files {
		log.Info(filename.Name())

		for index, missingFilename := range missingFiles {
			if filename.Name() == missingFilename {
				missingFiles = remove(missingFiles, index)
			}
		}
	}

	return len(missingFiles) == 0
}
