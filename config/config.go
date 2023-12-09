package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DebugMode bool
	LogFormat bool
	Port      string

	RequestTimeout int

	GithubApiBaseUrl  string
	GithubAccessToken string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to load .env file")
	}

	DebugMode = getBool("DEBUG_MODE", false)
	LogFormat = getBool("LOG_FORMAT", false)
	Port = getString("PORT", "8000")

	RequestTimeout = getInt("REQUEST_TIMEOUT", 60)

	GithubApiBaseUrl = getString("GITHUB_API_BASE_URL", "")
	GithubAccessToken = getString("GITHUB_ACCESS_TOKEN", "")
}

func getEnv(key string, fallback interface{}) interface{} {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getString(key string, fallback string) string {
	return getEnv(key, fallback).(string)
}

func getInt(key string, fallback int) int {
	f := strconv.Itoa(fallback)
	v, err := strconv.Atoi(getEnv(key, f).(string))
	if err != nil {
		return fallback
	}

	return v
}

func getBool(key string, fallback bool) bool {
	val := getEnv(key, strconv.FormatBool(fallback))
	valstring := val.(string)
	v, err := strconv.ParseBool(valstring)
	if err != nil {
		return fallback
	}

	return v
}
