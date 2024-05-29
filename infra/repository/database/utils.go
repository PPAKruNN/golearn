package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
)

const (
	ENV_FILE_PATH              = ".env"
	HOST_ENV                   = "DB_HOST"
	PORT_ENV                   = "DB_PORT"
	DATABASE_ENV               = "DB_DATABASE"
	USER_ENV                   = "DB_USER"
	PASSWORD_ENV               = "DB_PASSWORD"
	DATABASE_CONNECTION_STRING = "DB_CONN_STRING"
)

func loadDatabaseEnvs() pgx.ConnConfig {

	var config pgx.ConnConfig

	if os.Getenv("GO_ENV") == "prod" {
		fmt.Printf("\nEntering prod enviroment mode!\n")
		config = loadDatabaseEnvsUsingOS()
	} else {
		config = loadDatabaseEnvsUsingFile()
	}

	fmt.Printf("Using enviroment: %+v\n", config)

	return config

}

func loadDatabaseEnvsUsingFile() pgx.ConnConfig {

	env, err := godotenv.Read(ENV_FILE_PATH)
	if err != nil {
		panic(fmt.Sprintf("Could not read .ENV!! err: %v", err))
	}

	i, err := strconv.Atoi(env[PORT_ENV])

	if err != nil {
		panic(fmt.Sprintf("Could parse PORT ENV into a integer, err: %v", err))
	}

	return pgx.ConnConfig{
		Host:     env[HOST_ENV],
		Port:     uint16(i),
		Database: env[DATABASE_ENV],
		User:     env[USER_ENV],
		Password: env[PASSWORD_ENV],
	}
}

func loadDatabaseEnvsUsingOS() pgx.ConnConfig {

	err := godotenv.Load(ENV_FILE_PATH)

	dbconnstring := os.Getenv(DATABASE_CONNECTION_STRING)

	config, err := pgx.ParseConnectionString(dbconnstring)

	if err != nil {
		panic(fmt.Sprintf("Couldn't get enviroment variables!", err))
	}

	return config

}
