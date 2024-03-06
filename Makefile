
SHORT_ID := $(shell git rev-parse --short HEAD)
PACKAGES := $(shell go list ./...)
name := $(shell basename ${PWD})

export V_TAG = ${SHORT_ID}

# This assumes there is a .env file in the root of the project
define setup_env
  $(eval ENV_FILE := .env)
  $(eval include .env)
  $(eval export)
endef

all: help

## with-env: setup environment variables
.PHONY: with-env
with-env: 
	$(call setup_env)

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a make command to run"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## init: initialize project (make init module=github.com/user/project)
.PHONY: init
init:
	go mod init ${module}
	go install github.com/cosmtrek/air@latest

## vet: vet code
.PHONY: vet
vet:
	go vet $(PACKAGES)

# ## test: run unit tests
# .PHONY: test
# test:
# 	go test -race -cover $(PACKAGES)

# ## build: build a binary
# .PHONY: build
# build: test
# 	go build -o ./app -v

# ## docker-build: build project into a docker container image
# .PHONY: docker-build
# docker-build: test
# 	GOPROXY=direct docker buildx build -t ${name} .

# ## docker-run: run project in a container
# .PHONY: docker-run
# docker-run:
# 	docker run -it --rm -p 8080:8080 ${name}

# ## start: build and run local project
# .PHONY: start
# start: build
# 	air

## css: build tailwindcss
.PHONY: css
css:
	@npm run css

## css-dev: watch build tailwindcss
.PHONY: css-dev
css-dev:
	@npm run css:watch

## app-dev: run the app in development mode
.PHONY: app-dev
app-dev:
	@if command -v air > /dev/null; then \
	  air; \
	  echo "Watching...";\
	else \
	  read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	  if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	    go install github.com/cosmtrek/air@latest; \
	    air; \
	    echo "Watching...";\
	  else \
	    echo "You chose not to install air. Exiting..."; \
	    exit 1; \
	  fi; \
	fi

## run-dev: build and run the app in development mode
.PHONY: run-dev
make run-dev: 
	make -j 2 css-dev app-dev

## sqlc: generate Go code from SQL
.PHONY: sqlc
sqlc:
	@if command -v sqlc > /dev/null; then \
	  echo "Gerating Go code...";\
	  sqlc generate; \
	else \
	  read -p "Go's 'sqlc' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	  if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest; \
	  	echo "Gerating Go code...";\
	  	sqlc generate; \
	  else \
	    echo "You chose not to install sqlc. Exiting..."; \
	    exit 1; \
	  fi; \
	fi

## install-migrate: install migrate tool
.PHONY: install-migrate
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

## docker-up: start project in a container
.PHONY: docker-up
docker-up: with-env
	@docker compose up --build -d

## docker-down: stop project in a container
.PHONY: docker-down
docker-down: with-env
	@docker compose down

## docker-build: build project into a docker container image
.PHONY: docker-build
docker-build: with-env
	@docker build --target prod -t ${IMAGE_NAME}:${V_TAG} .

## docker-run: run project in a container
.PHONY: docker-run
docker-run: with-env
	@docker run --rm --name ${APP_NAME} -p 8080:8080 ${IMAGE_NAME}:${V_TAG}

## db-sh: open a psql shell within the postgres container
.PHONY: db-sh
db-sh: with-env
	@docker compose exec postgres psql -U ${DB_USER} -d ${DB_NAME}

## dump-schema: dump the database schema to a file db/schema.sql
.PHONY: dump-schema
dump-schema: with-env
	@docker compose exec postgres pg_dump -U postgres --schema-only -d ${DB_NAME} > db/schema.sql

# ## seed: seed the database
# .PHONY: seed
# seed: with-env
# 	@go run ./cmd/seed/main.go

## migrate: run migrations
.PHONY: migrate
migrate: with-env
	@if command -v migrate > /dev/null; then \
	  migrate -database ${DATABASE_URL} -path ./db/migrations up;\
	else \
	  echo "Installing 'migrate' tool...";\
	  make install-tools;\
	  echo "Running migrate up...";\
	  migrate -database ${DATABASE_URL} -path ./db/migrations up;\
	fi

## migrate-new: create a new migration. (make migrate-new name=create_users_table)
.PHONY: migrate-new
migrate-new:
	@if command -v migrate > /dev/null; then \
	  migrate create -ext sql -dir db/migrations -seq $(name);\
	else \
	  echo "Installing 'migrate' tool...";\
	  make install-migrate; \
	  echo "Creating new migration '$(name)'...";\
	  migrate create -ext sql -dir db/migrations -seq $(name);\
	fi

## migrate-down: run migrations down
.PHONY: migrate-down
migrate-down: with-env
	@migrate -database ${DATABASE_URL} -path ./db/migrations down

## migrate-drop: drop the database
.PHONY: migrate-drop
migrate-drop: with-env
	@migrate -database ${DATABASE_URL} -path ./db/migrations drop
