# komrifa

## Requirements

- Go 1.22
- Docker
- Tailwindcss

## Getting started

1. Clone the repo
2. Create the .env file
3. Run docker
4. Run migrations
5. Seed the database
6. Run the server

## Available commands

Run `make help` to get the list of commands

```
Choose a make command to run

  with-env          setup environment variables
  init              initialize project (make init module=github.com/user/project)
  vet               vet code
  css               build tailwindcss
  css-watch         watch build tailwindcss
  watch-api         watch for changes and restart the server
  sqlc              generate Go code from SQL
  install-migrate   install migrate tool
  docker-up         start project in a container
  docker-down       stop project in a container
  docker-build      build project into a docker container image
  docker-run        run project in a container
  db-sh             open a psql shell within the postgres container
  dump-schema       dump the database schema to a file db/schema.sql
  migrate           run migrations
  migrate-new       create a new migration. (make migrate-new name=create_users_table)
  migrate-down      run migrations down
  migrate-drop      drop the database
```
