FROM golang:1.19-alpine

WORKDIR /app

ADD go.mod .
ADD go.sum .

RUN ["go", "mod", "download"]

ADD main.go .

RUN ["go", "build", "."]

ENTRYPOINT ["/app/demoapp"]
