package main

import (
	"apiingolang/activity/bootconfig"
	"apiingolang/activity/business/utils/logging"
	"apiingolang/activity/db"
	"fmt"
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
	if os.Getenv("ENV") == "local" {
		if len(os.Getenv("PORT")) == 0 {
			os.Setenv("PORT", "9000")
		}
		if len(os.Getenv("CONFIGPATH")) == 0 {
			os.Setenv("CONFIGPATH", path.Join(pwd, "config/config.local.json"))
		}
		if len(os.Getenv("LOGDIR")) == 0 {
			os.Setenv("LOGDIR", path.Join(pwd, "logs"))
		}
	}
	fmt.Println("CONFIGPATH=", os.Getenv("CONFIGPATH"))
	fmt.Println("LOGDIR=", os.Getenv("LOGDIR"))
	bootconfig.InitConfig()
	// Loading DB connections
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}
	logging.NewLogger()
}
