# Собираем в гошке
FROM golang:1.15.2 as build

ENV BIN_FILE /opt/calendar/calendar-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}
RUN rm cmd/calendar/wire.go

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в apline образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -tags wireijnect \
        -o ${BIN_FILE} cmd/calendar/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="sterligov.denis94@yandex.ru"

ENV BIN_FILE "/opt/calendar/calendar-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/calendar/calendar_config.yml
COPY ./configs/calendar_config.yml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}