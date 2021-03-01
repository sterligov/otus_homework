FROM golang:1.15.2

WORKDIR "/migrations"

COPY ./migrations .

RUN go get -u github.com/pressly/goose/cmd/goose

CMD ["/go/bin/goose", "mysql", "calendar_user:calendar_pass@tcp(calendar_db:3306)/calendar?parseTime=true", "up"]
