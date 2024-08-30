### Use gfind if we are on macos
UNAME := $(shell uname)
ifeq ($(UNAME), Linux)
FIND_BIN=find
endif
ifeq ($(UNAME),Darwin)
FIND_BIN=gfind
endif

GOOS?=linux
OUTPUT_DIR?=./bin

GOPRIVATE=github.com/lrrountr
.EXPORT_ALL_VARIABLES:

.PHONY: fmt
fmt: $(SHELL $(FIND_BIN) -name "*.go")
		go fmt ./...

.PHONY: tidy
tidy:
		go mod tidy

.PHONY: build
build: fmt $(OUTPUT_DIR)/racing-draft

RACING_DRAFT_DB_PORT?=$$(docker port racing_draft_test_postgres 5234 | grep 0.0.0.0 | awk -F: '{print $$2}')
RACING_DRAFT_DB_HOST?=localhost
RACING_DRAFT_SKIP_SERVICES?=false

TEST_BASE_DIR = "./..."
GO_TEST_RUN?=

.PHONY: test
test: go-test

./bin/gotestsum:
		GOBIN=$(PWD)/bin go install gotest.tools/gotestsum@latest

.SHELLFLAGS = -ce

.PHONY: go-test
.ONESHELL:
go-test: ./bin/gotestsum ./bin/docker-compose
		workspace=$$(mktemp -d -p .)
		# We use a custom project name to prevent concurrent executions from clashing
		# Note we only generate lower case/numeric random strings, since docker/docker-compose
		# enforces that.
		project=racing-draft-$$(tr -dc a-z0-9 </dev/urandom | head -c 13 ; echo '')
		pg_container_name=$${project}_postgres_1
		if [ '-$(RACING_DRAFT_SKIP_SERVICES)' != '-true' ]; then
				cp ./cmd/server/tests/docker-compose.yaml $$workspace
				./bin/docker-compose -p $$project -f $$workspace/docker-compose.yaml up -d
				./scripts/wait_postgres.sh $${project}_postgres_1
				trap "./bin/docker-compose -p $$project -f $$workspace/docker-compose.yaml down --remove-orphans; rm -rf $$workspace" EXIT
		fi
		RACING_DRAFT_DB_HOST=$(RACING_DRAFT_DB_HOST) \
		RACING_DRAFT_DB_PORT=$$(docker port $$pg_container_name 5432 | grep 0.0.0.0 | awk -F: '{print $$2}') \
		RACING_DRAFT_DB_USER=postgres \
		RACING_DRAFT_DB_PASSWORD=postgres \
		RACING_DRAFT_TEST=1 \
		./bin/gotestsum -- -tags=swagger,test -run '$(GO_TEST_RUN)' $(TEST_BASE_DIR)

PLATFORM?=Linux
DOCKER_COMPOSE_VERSION=1.29.2

./bin/docker-compose:
		curl -L "https://github.com/docker/compose/releases/download/$(DOCKER_COMPOSE_VERSION)/docker-compose-$(PLATFORM)-x86_64" -o ./bin/docker-compose
		chmod +x ./bin/docker-compose