FROM golang:latest as base
FROM base as dev

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go mod download

EXPOSE 8000

RUN go build -o task_manager cmd/api/main.go
CMD [ "./taskmanager" ]