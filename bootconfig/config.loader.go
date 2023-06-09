package bootconfig

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

var ConfigManager IConfig

func InitConfig() {
	config, err := getFileConfigManager()
	if err != nil {
		log.Fatal("unable to load config")
	}
	ConfigManager = config
}

type fileConfig struct {
	keyVal map[string]interface{}
}

// loading configs from config file
func getFileConfigManager() (IConfig, error) {
	cfgMgr := fileConfig{}
	localConfig, err := os.Open(os.Getenv("CONFIGPATH"))
	if err != nil {
		return nil, err
	}
	rawbytes, err := ioutil.ReadAll(localConfig)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawbytes, &cfgMgr.keyVal)
	if err != nil {
		return nil, err
	}
	return &cfgMgr, nil
}

func (cfgmgr *fileConfig) Get(key string) ([]byte, error) {
	if val, ok := cfgmgr.keyVal[key]; ok {
		byteVal, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		return byteVal, nil
	}
	return nil, errors.New("no such config exists:\t" + key)
}

type IConfig interface {
	Get(key string) ([]byte, error)
}
