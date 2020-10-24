package util

import (
	"os"

	log "github.com/sirupsen/logrus"
)

//LogInfo logs an info level message with no additional fields
func LogInfo(m string) {
	log.WithFields(log.Fields{}).Info(m)
}

//LogError logs a message with the provided error string as a field
func LogError(m string, e string) {
	log.WithFields(log.Fields{
		"Error": e,
	}).Error(m)
}

//LogDebug only runs when the env var `DEBUG` is set to "true"
// It adds the "Debug" field and logs an info level message
// This is a hack to avoid setting up a full-blown logger for now, eventually we should do that
func LogDebug(m string) {
	if os.Getenv("DEBUG") == "true" {
		log.WithFields(log.Fields{
			"Debug": true,
		}).Info(m)
	}
}
