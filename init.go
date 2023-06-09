package main

import (
	"apiingolang/activity/bootconfig"
	"apiingolang/activity/db"
	"fmt"
	"os"
)

const (
	app = "activity-selector"
)

func init() {
	os.Setenv("APP", app)
	os.Setenv("PORT", "9000")
	os.Setenv("CONFIGPATH", "/Users/rahulsahani/Desktop/Codebase/repos/apiingolang/config/config.local.json")

	bootconfig.InitConfig()
	// Loading DB connections
	if err := db.Init(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
