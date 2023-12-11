# Cube dashboard backend

The backend provides authentication and data ingestion handling for the cube dashboard.

## Run locally

Make sure you have the `.env` file in project root and run `make run`. The ENVs are auto loaded by the `dotenv` package.

## Build

Make sure you have the `.env` file in project root and run `make build`. The ENVs are auto loaded by the `dotenv` package.

ENVs:

```sh
POSTGRES_HOST=your_postgres_host
POSTGRES_PORT=your_postgres_port
POSTGRES_USER=your_postgres_user
POSTGRES_PASSWORD=your_postgres_password
POSTGRES_DB_NAME=your_postgres_db_name

PORT=your_port

GOOGLE_MAP_API_KEY=your_google_map_api_key
CUBE_API_SECRET=cube_api_secret
API_JWT_SECRET=api_jwt_secret
```
