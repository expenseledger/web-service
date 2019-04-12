# Expense Ledger Web Service

The backend server for Expense Ledger application. It was developed originally
to solve the pain of the lack of good-enough expense tracking applications.

## Prerequisites
* [Go](https://golang.org/)
* [Dep](https://github.com/golang/dep)
* [Docker Compose](https://docs.docker.com/compose/) (_This is just for
local database_)

## Running

1. Start the containers

```bash
docker-compose up -d
```

This will start a PostgreSQL container and mount the data in `./pgdata/`
directory. It also start pgadmin4 container, just in case you need a client
to manage the database.

2. Initial the database

```bash
# this will be added later
# for now, just create a database named `expense_ledger_web_service`
```

3. Copy and edit `.env` file

```bash
cp .env.sample .env
vim .env # edit the .env file
```

4. Run the main process

```bash
go run main.go
```
