package configs

import "os"

var (
	Port             string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDBName   string
)

func InitEnvConfig() {
	Port = os.Getenv("PORT")

	if Port == "" {
		panic("PORT env variable is not set")
	}

	PostgresHost = os.Getenv("POSTGRES_HOST")
	PostgresPort = os.Getenv("POSTGRES_PORT")
	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	PostgresDBName = os.Getenv("POSTGRES_DB_NAME")

	if PostgresHost == "" {
		panic("POSTGRES_HOST env variable is not set")
	}

	if PostgresPort == "" {
		panic("POSTGRES_PORT env variable is not set")
	}

	if PostgresUser == "" {
		panic("POSTGRES_USER env variable is not set")
	}

	if PostgresPassword == "" {
		panic("POSTGRES_PASSWORD env variable is not set")
	}

	if PostgresDBName == "" {
		panic("POSTGRES_DB_NAME env variable is not set")
	}

}
