package utility

import (
	"os"
)

var DatabasePath = "test.sqlite3"
var DebugEnvironment = true

func InitConfig() {
	databasePath := os.Getenv("LOGI_TRACKER_DATABASE_PATH")
	if len(databasePath) > 0 {
		DatabasePath = databasePath
		DebugEnvironment = false
	}
}
