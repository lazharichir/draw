package config

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName = "draw" // change to relevant project name

type Config struct {
	R2_ACCOUNT_ID                 string
	R2_ACCESS_KEY_ID              string
	R2_ACCESS_KEY_SECRET          string
	R2_TILECACHE_BUCKET_NAME      string
	R2_TILECACHE_BUCKET_NAME_TEST string
}

var config *Config

func init() {
	config = getConfig()
}

func loadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func getConfig() *Config {
	loadEnv()

	R2_ACCOUNT_ID := os.Getenv("R2_ACCOUNT_ID")
	R2_ACCESS_KEY_ID := os.Getenv("R2_ACCESS_KEY_ID")
	R2_ACCESS_KEY_SECRET := os.Getenv("R2_ACCESS_KEY_SECRET")
	R2_TILECACHE_BUCKET_NAME := os.Getenv("R2_TILECACHE_BUCKET_NAME")
	R2_TILECACHE_BUCKET_NAME_TEST := os.Getenv("R2_TILECACHE_BUCKET_NAME_TEST")

	return &Config{
		R2_ACCOUNT_ID:                 R2_ACCOUNT_ID,
		R2_ACCESS_KEY_ID:              R2_ACCESS_KEY_ID,
		R2_ACCESS_KEY_SECRET:          R2_ACCESS_KEY_SECRET,
		R2_TILECACHE_BUCKET_NAME:      R2_TILECACHE_BUCKET_NAME,
		R2_TILECACHE_BUCKET_NAME_TEST: R2_TILECACHE_BUCKET_NAME_TEST,
	}
}

func GetConfig() *Config {
	return config
}
