package web

import (
	log "github.com/sirupsen/logrus"
)

// Remove element from string slice, order is not preserved.
func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// FilesExist returns false if the minimum files required for web frontend delivery
// have not been compiled into the web binary.
func FilesExist() bool {
	missingFiles := []string{"index.html.tmpl", "app.js", "app.css"}

	files, err := ContentFS.ReadDir("static")
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
