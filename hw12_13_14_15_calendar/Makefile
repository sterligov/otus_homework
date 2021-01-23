include ./deployments/.env

CALENDAR_BIN := "./bin/calendar"
SCHEDULER_BIN := "./bin/scheduler"
SENDER_BIN := "./bin/sender"
CALENDAR_DOCKER_IMG="calendar:develop"
SCHEDULER_DOCKER_IMG="scheduler:develop"
SENDER_DOCKER_IMG="sender:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build: build-calendar build-scheduler build-sender

build-calendar:
	go build -v -o $(CALENDAR_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-scheduler:
	go build -v -o $(SCHEDULER_BIN) -ldflags "$(LDFLAGS)" ./cmd/scheduler

build-sender:
	go build -v -o $(SENDER_BIN) -ldflags "$(LDFLAGS)" ./cmd/sender

run-calendar: build-calendar
	$(CALENDAR_BIN) -config ./configs/calendar_config.yml

run-scheduler: build-scheduler
	$(SCHEDULER_BIN) -config ./configs/scheduler_config.yml

run-sender: build-sender
	$(SENDER_BIN) -config ./configs/sender_config.yml

run: run-scheduler run-calendar run-sender

build-calendar-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(CALENDAR_DOCKER_IMG) \
		-f build/Calendar.Dockerfile .

build-scheduler-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(SCHEDULER_DOCKER_IMG) \
		-f build/Scheduler.Dockerfile .

build-sender-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(SENDER_DOCKER_IMG) \
		-f build/Sender.Dockerfile .

run-calendar-img: build-calendar-img
	docker run $(CALENDAR_DOCKER_IMG)

run-scheduler-img: build-scheduler-img
	docker run $(SCHEDULER_DOCKER_IMG)

run-sender-img: build-sender-img
	docker run $(SENDER_DOCKER_IMG)

calendar-version: build
	$(CALENDAR_BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.30.0

lint: install-lint-deps
	golangci-lint run ./...

wire:
	wire cmd/calendar/wire.go
	wire cmd/scheduler/wire.go
	wire cmd/sender/wire.go

docker-up:
	cd deployments && docker-compose up -d

docker-stop:
	cd deployments && docker-compose stop

migrations:
	goose -dir migrations mysql "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=true" up

generate:
	go generate ./...

.PHONY: migrations build run build-img run-img version test lint
