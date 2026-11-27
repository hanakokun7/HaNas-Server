package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"text/template"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AppConfig struct {
	BACKEND_URL string `json:"BACKEND_URL"`
	Db_type     string `json:"db_type"`
	DbFile      string `json:"dbFile"`
	Db_url      string `json:"db_url"`
	DataDir     string `json:"dataDir"`
	ThumbDir    string `json:"thumbDir"`
	Https_f     bool   `json:"https_f"`
	Cert        string `json:"cert"`
	Key         string `json:"key"`
}

var acfg AppConfig

func Init() string {
	var err error
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}
	dir := filepath.Dir(exePath)
	configPath := filepath.Join(dir, "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}
	if err := json.Unmarshal(data, &acfg); err != nil {
		return ""
	}
	switch acfg.Db_type {
	case "mysql":
		db, err = gorm.Open(mysql.Open(acfg.Db_url), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	default:
		db, err = gorm.Open(sqlite.Open(acfg.DbFile), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	}
	db.AutoMigrate(&Config{}, &User{}, &Node{}, &Share{}, &Trash{})
	if err := initJWTSecret(); err != nil {
		panic(err)
	}
	tmpl, err := template.New("config.js").Parse(configJS)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, acfg); err != nil {
		return ""
	}
	js_script := buf.String()
	return js_script
}
