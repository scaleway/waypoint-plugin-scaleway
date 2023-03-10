FROM golang:1.19-alpine3.15 as build

# Install the Protocol Buffers compiler and Go plugin
RUN apk add protobuf git make zip
RUN go install github.com/golang/protobuf/protoc-gen-go@v1.5.2 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

# Create the source folder
RUN mkdir /go/plugin
WORKDIR /go/plugin

# Download dependencies
COPY go.mod /go/plugin
COPY go.sum /go/plugin
RUN go mod download

# Copy the source to the build folder
COPY . /go/plugin

# Build the plugin
ARG VERSION
ENV VERSION=$VERSION
RUN chmod +x ./print_arch
RUN make all

# Create the zipped binaries
RUN make zip

FROM scratch as export_stage

COPY --from=build /go/plugin/bin/*.zip ./
