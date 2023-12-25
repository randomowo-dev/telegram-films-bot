package config

var (
	KinopoiskApiUnofficialUrl = getEnv(
		"KINOPOISK_API_UNOFFICIAL_URL",
		"https://kinopoiskapiunofficial.tech",
	)
	KinopoiskApiUnofficialToken           = getEnv("KINOPOISK_API_UNOFFICIAL_TOKEN", "")
	KinopoiskApiUnofficialFilmsApiVersion = getEnv("KINOPOISK_API_UNOFFICIAL_FILMS_API_VERSION", "v2.2")
	KinopoiskApiUnofficialStaffApiVersion = getEnv("KINOPOISK_API_UNOFFICIAL_STAFF_API_VERSION", "v1")
)
