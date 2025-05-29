package config

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
)

type AppConfig struct {
	version    string
	serverPort string
}

func ReadAppConfig() (*AppConfig, error) {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return nil, fmt.Errorf("config/appConfig.ReadAppConfig: [Didn't manage to read build info]")
	}

	serverPort, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		return nil, fmt.Errorf("config/appConfig.ReadAppConfig: [SERVER_PORT is not defined]")
	}
	if _, err := strconv.ParseUint(serverPort, 10, 16); err != nil {
		return nil, fmt.Errorf("config/appConfig.ReadAppConfig: [SERVER_PORT is should be uint16]")
	}

	return &AppConfig{
		version:    buildInfo.Main.Version,
		serverPort: serverPort,
	}, nil
}

func (ac *AppConfig) Version() string {
	return ac.version
}

func (ac *AppConfig) ServerPort() string {
	return ac.serverPort
}
