package config

import (
	"fmt"
	"strconv"
)

var (
	serverPort, _       = strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	ServerAddr          = fmt.Sprintf(":%d", serverPort)
	ServerAuthSecret    = []byte(getEnv("SERVER_AUTH_SECRET", ""))
	ServerBasicUsername = getEnv("SERVER_AUTH_USERNAME", "")
	ServerBasicPassword = getEnv("SERVER_AUTH_PASSWORD", "")
)
