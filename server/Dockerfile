FROM golang:latest as base
FROM base as dev

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go mod download

RUN go build -o task_manager .
CMD [ "./task_manager" ]