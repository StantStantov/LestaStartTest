package config

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
)

const (
	ServerPortEnv  = "SERVER_PORT"
	DatabaseUrlEnv = "DATABASE_URL"
	PathToDocsEnv  = "DOCUMENTS_PATH"
)

type AppConfig struct {
	version    string
	serverPort string
	dbUrl      string
	pathToDocs string
}

func (c *AppConfig) Version() string {
	return c.version
}

func (c *AppConfig) ServerPort() string {
	return c.serverPort
}

func (c *AppConfig) DatabaseUrl() string {
	return c.dbUrl
}

func (c *AppConfig) PathToDocuments() string {
	return c.pathToDocs
}

func ReadAppConfig() (*AppConfig, error) {
	buildInfo, err := readBuildInfo()
	if err != nil {
		return nil, fmt.Errorf("config/appConfig.ReadAppConfig: [Didn't manage to read build info]")
	}

	serverPort, err := readServerPort()
	if err != nil {
		return nil, fmt.Errorf("config/appConfig.ReadAppConfig: [%w]", err)
	}

	dbUrl, err := readDatabaseUrl()
	if err != nil {
		return nil, fmt.Errorf("config/appConfig.ReadAppConfig: [%w]", err)
	}

	pathToDocs, err := readPathToDocuments()
	if err != nil {
		return nil, fmt.Errorf("config/appConfig.ReadAppConfig: [%w]", err)
	}

	return &AppConfig{
		version:    buildInfo.Main.Version,
		serverPort: serverPort,
		dbUrl:      dbUrl,
		pathToDocs: pathToDocs,
	}, nil
}

func readBuildInfo() (*debug.BuildInfo, error) {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok || buildInfo == nil {
		return nil, fmt.Errorf("config/appConfig.ReadAppConfig: [Didn't manage to read build info]")
	}

	return buildInfo, nil
}

func readServerPort() (string, error) {
	serverPort, ok := os.LookupEnv(ServerPortEnv)
	if !ok {
		return "", fmt.Errorf("config/appConfig.readServerPort: [ENV %s is not defined]", serverPort)
	}
	if _, err := strconv.ParseUint(serverPort, 10, 16); err != nil {
		return "", fmt.Errorf("config/appConfig.readServerPort: [ENV %s is should be uint16]", serverPort)
	}

	return serverPort, nil
}

func readDatabaseUrl() (string, error) {
	dbUrl, ok := os.LookupEnv(DatabaseUrlEnv)
	if !ok {
		return "", fmt.Errorf("config/appConfig.readDatabaseUrl: [ENV %s is not defined]", DatabaseUrlEnv)
	}

	return dbUrl, nil
}

func readPathToDocuments() (string, error) {
	pathToDocs, ok := os.LookupEnv(PathToDocsEnv)
	if !ok {
		return "", fmt.Errorf("config/appConfig.readDatabaseUrl: [ENV %s is not defined]", PathToDocsEnv)
	}

	return pathToDocs, nil
}
