package config

import (
	"os"
)

var (
	// kinopoiskapiunofficial.tech
	KinopoiskApiUnofficialUrl = getEnv(
		"KINOPOISK_API_UNOFFICIAL_URL",
		"https://kinopoiskapiunofficial.tech",
	)
	KinopoiskApiUnofficialToken           = getEnv("KINOPOISK_API_UNOFFICIAL_TOKEN", "")
	KinopoiskApiUnofficialFilmsApiVersion = getEnv("KINOPOISK_API_UNOFFICIAL_FILMS_API_VERSION", "v2.2")
	KinopoiskApiUnofficialStaffApiVersion = getEnv("KINOPOISK_API_UNOFFICIAL_STAFF_API_VERSION", "v1")
)

func getEnv(key, defaultValue string) string {
	if value, found := os.LookupEnv(key); found {
		return value
	}
	return defaultValue
}
