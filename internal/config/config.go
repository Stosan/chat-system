package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"github.com/joho/godotenv"
)

// Config stores the application configuration from environment variables
type ConfigApplication struct {
	Env                 string
	APIKey              string
	PORT                string
	DBName              string
	DBConnURL           string
	DBHost              string
	DBUsername          string
	DBPassword          string
	DBSSLMode           string
	CloudinaryCloudName string 
	CloudinaryAPIKey    string 
	CloudinaryAPISecret string 
	RedisAddress        string // Added field for Redis
	RedisUsername       string // Added field for Redis
	RedisPassword       string // Added field for Redis
}

var AppConfig ConfigApplication

// findRootDir finds the project root directory by looking for go.mod
func findRootDir() (string, error) {
	// Get the directory of the current file
	_, filename, _, _ := runtime.Caller(0)
	currentDir := path.Dir(filename)

	// Walk up until we find go.mod
	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return "", fmt.Errorf("could not find project root (no go.mod found)")
		}
		currentDir = parentDir
	}
}

// init function is used to initialize the application configuration.
// It loads environment variables from .env files located at different paths.
// The function iterates over the paths and loads the first .env file it finds.
// If it fails to load an .env file from any of the paths, it logs a fatal error.
// After successfully loading the .env file, it assigns the environment variables to the AppConfig struct.
func init() {
	// Find project root and load .env
	rootDir, err := findRootDir()
	if err != nil {
		log.Fatal(err)
	}

	var (
		present bool
	)
	err = godotenv.Load(filepath.Join(rootDir, ".env"))
	if err != nil {
		log.Fatal(err)
	}

	AppConfig.Env, present = os.LookupEnv("GO_ENV")
	if !present {
		panic("GO_ENV environment variable is not set")
	}

	AppConfig.PORT, present = os.LookupEnv("GO_PORT")
	if !present {
		panic("GO_PORT environment variable is not set")
	}

	AppConfig.APIKey, present = os.LookupEnv("API_KEY")
	if !present {
		panic("API_KEY environment variable is not set")
	}

	AppConfig.DBName, present = os.LookupEnv("DB_NAME")
	if !present {
		panic("DB_NAME environment variable is not set")
	}

	AppConfig.DBHost, present = os.LookupEnv("DB_HOST")
	if !present {
		panic("DB_HOST environment variable is not set")
	}

	AppConfig.DBUsername, present = os.LookupEnv("DB_USERNAME")
	if !present {
		panic("DB_USERNAME environment variable is not set")
	}

	AppConfig.DBPassword, present = os.LookupEnv("DB_PASSWORD")
	if !present {
		panic("DB_PASSWORD environment variable is not set")
	}

	AppConfig.DBSSLMode, present = os.LookupEnv("DB_SSLMODE")
	if !present {
		panic("DB_SSLMODE environment variable is not set")
	}

	AppConfig.CloudinaryCloudName, present = os.LookupEnv("CLOUDINARY_CLOUD_NAME") // Added for Cloudinary
	if !present {
		panic("CLOUDINARY_CLOUD_NAME environment variable is not set")
	}

	AppConfig.CloudinaryAPIKey, present = os.LookupEnv("CLOUDINARY_API_KEY") // Added for Cloudinary
	if !present {
		panic("CLOUDINARY_API_KEY environment variable is not set")
	}

	AppConfig.CloudinaryAPISecret, present = os.LookupEnv("CLOUDINARY_API_SECRET") // Added for Cloudinary
	if !present {
		panic("CLOUDINARY_API_SECRET environment variable is not set")
	}

	AppConfig.RedisAddress, present = os.LookupEnv("REDIS_ADDRESS") // Added for Redis
	if !present {
		panic("REDIS_ADDRESS environment variable is not set")
	}

	AppConfig.RedisUsername, present = os.LookupEnv("REDIS_USERNAME") // Added for Redis
	if !present {
		panic("REDIS_USERNAME environment variable is not set")
	}

	AppConfig.RedisPassword, present = os.LookupEnv("REDIS_PASSWORD") // Added for Redis
	if !present {
		panic("REDIS_PASSWORD environment variable is not set")
	}

}
