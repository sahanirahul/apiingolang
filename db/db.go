package db

import (
	"apiingolang/activity/bootconfig"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func Init() error {
	secretString, err := bootconfig.ConfigManager.Get("db/postgres")
	if err != nil || secretString == nil {
		if err == nil {
			err = errors.New("null_secretString")
		}
		return err
	}
	err = initPostgres(secretString)
	if err != nil {
		return err
	}
	return nil
}

const (
	maxOpenConnection    = 10
	maxIdleConnection    = 5
	maxConnectionTimeout = time.Hour * 1
)

var ClientActivity *sql.DB

type DbObject struct {
	User                 string `json:"username" validate:"required"`
	Password             string `json:"password" validate:"required"`
	Engine               string `json:"engine" validate:"required"`
	Host                 string `json:"host" validate:"required"`
	Port                 int    `json:"port" validate:"required"`
	DbInstanceIdentifier string `json:"dbInstanceIdentifier" validate:"required"`
	DbName               string `json:"dbName" validate:"required"`
}

func initPostgres(secret []byte) error {
	fmt.Println("init_db_start")
	var dbObject DbObject
	err := json.Unmarshal(secret, &dbObject)
	if err != nil {
		fmt.Println("Not_able_to_structure_db_secretString" + err.Error())
		return err
	}
	ClientActivity, err = getPostgresConnection(dbObject)
	if err != nil {
		return err
	}
	return nil
}

func getPostgresConnection(dbObject DbObject) (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbObject.Host, dbObject.Port, dbObject.User, dbObject.Password, dbObject.DbName)

	Client, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = Client.Ping()
	if err != nil {
		return nil, err
	}

	Client.SetMaxOpenConns(maxOpenConnection)
	Client.SetMaxIdleConns(maxIdleConnection)
	Client.SetConnMaxLifetime(maxConnectionTimeout)
	if err != nil {
		return nil, err
	}
	fmt.Println("PostgreSQL Connection established with: ", psqlInfo)

	return Client, nil
}
