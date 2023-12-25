package config

import (
	"fmt"
	"strconv"
)

var (
	dbUsername = getEnv("DB_USERNAME", "")
	dbPassword = getEnv("DB_PASSWORD", "")
	dbHost     = getEnv("DB_HOST", "")
	dbPort, _  = strconv.Atoi(getEnv("DB_PORT", "27017"))
	DbUrl      = fmt.Sprintf("mongodb://%s:%s@%s:%d", dbUsername, dbPassword, dbHost, dbPort)
	DbName     = getEnv("DB_NAME", "")
)
