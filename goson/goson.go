package goson

/*
	TODO: Better error handling
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

const (
	defaultConfigPath     = ".config"
	defaultConfigFileName = "config.json"
)

type Config struct {
	AppName  string
	Loaded   bool
	Saved    bool
	Contents map[string]interface{}
}

// Get the filepath for the config file
func (conf Config) FilePath() (fpath string, err error) {
	u, err := user.Current()
	if u != nil && err == nil {
		fpath = filepath.Join(u.HomeDir, defaultConfigPath, conf.AppName, defaultConfigFileName)
	}
	return
}

// Load configuration from filesystem
func (conf *Config) Load() {
	conf.Contents = make(map[string]interface{})

	confFilePath, err := conf.FilePath()
	fileBody, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		fmt.Println(err)

	}

	err = json.Unmarshal(fileBody, &conf.Contents)
	if err != nil {
		fmt.Println(err)
	}

	conf.Loaded = true
	conf.Saved = true
}

// Retrieves Json data
func (conf Config) JsonData() ([]byte, error) {
	return json.MarshalIndent(conf.Contents, "", "	")
}

// Write config file back to filesystem
func (conf Config) Write() {
	data, err := conf.JsonData()
	confFilePath, err := conf.FilePath()
	if err != nil {
		fmt.Println(err)
		return
	}

	confDirname := filepath.Dir(confFilePath)
	err = os.MkdirAll(confDirname, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(confFilePath, data, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	conf.Saved = true
}

func (conf *Config) Clear() {
	conf.Contents = make(map[string]interface{})
	conf.Saved = false
}

// Gets a value from the config
func (conf Config) Get(key string) (interface{}, bool) {
	v, ok := conf.Contents[key]
	return v, ok
}

// Gets a value from the config
func (conf Config) GetString(key string) (string, bool) {
	v, ok := conf.Contents[key].(string)
	return v, ok
}

// Gets a value from the config
func (conf Config) GetInt(key string) (int, bool) {
	v, ok := conf.Contents[key].(int)
	return v, ok
}

// Returns a list
func (conf Config) GetList(key string) (l []interface{}, ok bool) {
	v, ok := conf.Get(key)
	if ok {
		l, ok = v.([]interface{})
	}

	return l, ok
}

// Sets a value in the config
func (conf *Config) Set(key string, value interface{}) {
	conf.Contents[key] = value
	conf.Saved = false
}

// Sets value and writes to file
func (conf *Config) SetAndWrite(key string, value interface{}) {
	conf.Set(key, value)
	conf.Write()
}

// Load config for app name
func LoadConfig(appName string) (conf Config) {
	conf.AppName = appName
	conf.Load()
	return conf
}
