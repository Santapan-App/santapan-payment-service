# Database
POSTGRES_USER ?= user
POSTGRES_PASSWORD ?= s4nt4p4nDatab4s3
POSTGRES_ADDRESS ?= localhost:5432
POSTGRES_DATABASE ?= santapan_db

# Exporting bin folder to the path for makefile
export PATH   := $(PWD)/bin:$(PATH)
# Default Shell
export SHELL  := bash
# Type of OS: Linux or Darwin.
export OSTYPE := $(shell uname -s | tr A-Z a-z)
export ARCH := $(shell uname -m)



# --- Tooling & Variables ----------------------------------------------------------------
include ./misc/make/tools.Makefile
include ./misc/make/help.Makefile

# ~~~ Development Environment ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

up: dev-env dev-air             ## Startup / Spinup Docker Compose and air
down: docker-stop               ## Stop Docker
destroy: docker-teardown clean  ## Teardown (removes volumes, tmp files, etc...)

install-deps: migrate air gotestsum tparse mockery ## Install Development Dependencies (localy).
deps: $(MIGRATE) $(AIR) $(GOTESTSUM) $(TPARSE) $(MOCKERY) $(GOLANGCI) ## Checks for Global Development Dependencies.
deps:
	@echo "Required Tools Are Available"

dev-env: ## Bootstrap Environment (with a Docker-Compose help).
	@ docker compose up -d --build postgres

dev-env-test: dev-env ## Run application (within a Docker-Compose help)
	@ $(MAKE) image-build
	docker compose up web

dev-air: $(AIR) ## Starts AIR ( Continuous Development app).
	air

docker-stop:
	@ docker compose down

docker-teardown:
	@ docker compose down --remove-orphans -v

# ~~~ Code Actions ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

lint: $(GOLANGCI) ## Runs golangci-lint with predefined configuration
	@echo "Applying linter"
	golangci-lint version
	golangci-lint run -c .golangci.yaml ./...

# -trimpath - will remove the filepathes from the reports, good to same money on network trafic,
#             focus on bug reports, and find issues fast.
# - race    - adds a racedetector, in case of racecondition, you can catch report with sentry.
#             https://golang.org/doc/articles/race_detector.html
#
# todo(butuzov): add additional flags to compiler to have an `version` flag.
build: ## Builds binary
	@ printf "Building aplication... "
	@ go build \
		-trimpath  \
		-o engine \
		./app/
	@ echo "done"


build-race: ## Builds binary (with -race flag)
	@ printf "Building aplication with race flag... "
	@ go build \
		-trimpath  \
		-race      \
		-o engine \
		./app/
	@ echo "done"


go-generate: $(MOCKERY) ## Runs go generte ./...
	go generate ./...


TESTS_ARGS := --format testname --jsonfile gotestsum.json.out
TESTS_ARGS += --max-fails 2
TESTS_ARGS += -- ./...
TESTS_ARGS += -test.parallel 2
TESTS_ARGS += -test.count    1
TESTS_ARGS += -test.failfast
TESTS_ARGS += -test.coverprofile   coverage.out
TESTS_ARGS += -test.timeout        5s
TESTS_ARGS += -race

tests: $(GOTESTSUM)
	@ gotestsum $(TESTS_ARGS) -short

tests-complete: tests $(TPARSE) ## Run Tests & parse details
	@cat gotestsum.json.out | $(TPARSE) -all -notests

# ~~~ Docker Build ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.ONESHELL:
image-build:
	@ echo "Docker Build"
	@ DOCKER_BUILDKIT=0 docker build \
		--file Dockerfile \
		--tag koobam-user \
			.

# Commenting this as this not relevant for the project, we load the DB data from the SQL file.
# please refer this when introducing the database schema migrations.

# # ~~~ Database Migrations ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

POSTGRES_DSN := "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DATABASE)?sslmode=disable"

migrate-up: ## Apply all (or N) migrations.
	@ echo "Applying migrations..."
	@ migrate -path=migrations -database ${POSTGRES_DSN} up

migrate-down: ## Roll back all (or N) migrations.
	@ echo "Rolling back migrations..."
	@ migrate -path=migrations -database ${POSTGRES_DSN} down

migrate-drop: ## Drop all data and schema.
	@ echo "Dropping all database schema..."
	@ migrate -path=migrations -database ${POSTGRES_DSN} drop

migrate-create: ## Create a new migration.
	@ read -p "Please provide a name for the migration: " Name; \
	migrate create -ext sql -dir migrations -seq $${Name}

# MYSQL_DSN := "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)"

# migrate-up: $(MIGRATE) ## Apply all (or N up) migrations.
# 	@ read -p "How many migration you wants to perform (default value: [all]): " N; \
# 	migrate  -database $(MYSQL_DSN) -path=misc/migrations up ${NN}

# .PHONY: migrate-down
# migrate-down: $(MIGRATE) ## Apply all (or N down) migrations.
# 	@ read -p "How many migration you wants to perform (default value: [all]): " N; \
# 	migrate  -database $(MYSQL_DSN) -path=misc/migrations down ${NN}

# .PHONY: migrate-drop
# migrate-drop: $(MIGRATE) ## Drop everything inside the database.
# 	migrate  -database $(MYSQL_DSN) -path=misc/migrations drop

# .PHONY: migrate-create
# migrate-create: $(MIGRATE) ## Create a set of up/down migrations with a specified name.
# 	@ read -p "Please provide name for the migration: " Name; \
# 	migrate create -ext sql -dir misc/migrations $${Name}

# ~~~ Cleans ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

clean: clean-artifacts clean-docker

clean-artifacts: ## Removes Artifacts (*.out)
	@printf "Cleanning artifacts... "
	@rm -f *.out
	@echo "done."

clean-docker: ## Removes dangling docker images
	@ docker image prune -f

# Path to protobuf files and generated code
PROTO_FILES := $(wildcard internal/proto/*.proto)
PROTO_GEN_PATH := .

generate-proto: ## Generate Go code from protobuf definitions
	@ mkdir -p $(PROTO_GEN_PATH)
	@ protoc \
		--go_out=$(PROTO_GEN_PATH) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GEN_PATH) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@ echo "Protobuf code generated in $(PROTO_GEN_PATH)"