package utility

import (
	"os"
)

var UrlRoot = "http://127.0.0.1"
var DatabasePath = "test.sqlite3"
var DebugEnvironment = true

func InitConfig() {
	urlRoot := os.Getenv("LOGI_TRACKER_URLROOT")
	databasePath := os.Getenv("LOGI_TRACKER_URLROOT")

	if len(urlRoot) > 0 {
		UrlRoot = urlRoot
		DebugEnvironment = false
	}
	if len(databasePath) > 0 {
		DatabasePath = databasePath
		DebugEnvironment = false
	}
}