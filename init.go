package main

import (
	"apiingolang/activity/bootconfig"
	"apiingolang/activity/business/utils/logging"
	"apiingolang/activity/db"
	"log"
	"os"
	"path"
)

const (
	app = "activity-selector"
)

func init() {
	pwd, _ := os.Getwd()
	os.Setenv("APP", app)
	os.Setenv("PORT", "9000")
	os.Setenv("CONFIGPATH", "/Users/rahulsahani/Desktop/Codebase/repos/apiingolang/config/config.local.json")
	os.Setenv("LOGPATH", path.Join(pwd, "/logs/activity.log"))

	bootconfig.InitConfig()
	// Loading DB connections
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}
	logging.NewLogger()
}
