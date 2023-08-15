package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db  *gorm.DB
	cfg Config
)

type Config struct {
	Env        string     `yaml:"env" env:"ENV" env-default:"local"`
	HTTPServer HTTPServer `yaml:"http_server"`
	DataBase   DataBase   `yaml:"db"`
}
type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DataBase struct {
	DbUser     string `yaml:"db_user" env-required:"true"`
	DbPassword string `yaml:"db_password"`
	DbHost     string `yaml:"db_host" env-default:"localhost:3306"`
}

func Connect() {
	// Checking if the connection is open
	// sqlDB := db.DB()
	// if sqlDB != nil {
	// 	return
	// }
	cfg := GetCFG()
	dbConfig := fmt.Sprintf("%s:%s@tcp(%s)/test?parseTime=true",
		cfg.DataBase.DbUser, cfg.DataBase.DbPassword, cfg.DataBase.DbHost)
	d, err := gorm.Open("mysql", dbConfig)
	if err != nil {
		panic(err)
	}
	db = d
}
func MustLoad() *Config {
	if cfg.IsInitialized() {
		return &cfg
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file %s does not exist", configPath)
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config file %s: %s", configPath, err)
	}

	return &cfg
}

func GetDB() *gorm.DB {
	return db
}

func GetCFG() *Config {
	return &cfg
}
func (cfg Config) IsInitialized() bool {
	return cfg.DataBase.DbUser != ""
}
