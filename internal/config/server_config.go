package config

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

var (
	serverPort, _        = strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	ServerAddr           = fmt.Sprintf(":%d", serverPort)
	ServerAuthSecret     = []byte(getEnv("SERVER_AUTH_SECRET", ""))
	serverBasicUsername  = getEnv("SERVER_AUTH_USERNAME", "")
	serverBasicPassword  = getEnv("SERVER_AUTH_PASSWORD", "")
	ServerBasicAuthToken = base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf(
			"%s:%s",
			serverBasicUsername,
			serverBasicPassword,
		)),
	)
)
