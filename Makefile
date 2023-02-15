PLUGIN_NAME=scaleway

ifndef _ARCH
_ARCH := $(shell ./print_arch)
export _ARCH
endif

.PHONY: all

# Plugin Version
VERSION?=

LDFLAGS=-X main.Version=${VERSION}

PLUGINS=\
	container


all: protos build

protos: ${PLUGINS:%=protos-%}

protos-%: %
	@echo ""
	@echo Build Protos $<

	protoc -I thirdparty/proto -I . --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./$</plugin.proto

build: ${PLUGINS:%=build-%}

build-%: %
	@echo ""
	@echo Compile Plugin $<

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o ./bin/linux_amd64/waypoint-plugin-$(PLUGIN_NAME)-$< ./cmd/$</main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o ./bin/darwin_amd64/waypoint-plugin-$(PLUGIN_NAME)-$< ./cmd/$</main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o ./bin/windows_amd64/waypoint-plugin-$(PLUGIN_NAME)-$<.exe ./cmd/$</main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="${LDFLAGS}" -o ./bin/windows_386/waypoint-plugin-$(PLUGIN_NAME)-$<.exe ./cmd/$</main.go

install: ${PLUGINS:%=install-%}

# Install the plugin locally
install-%: %
	@echo ""
	@echo "Installing Plugin"

	cp ./bin/${_ARCH}_amd64/waypoint-plugin-${PLUGIN_NAME}-$< ${HOME}/.config/waypoint/plugins/

zip: ${PLUGINS:%=zip-%}

# Zip the built plugin binaries
zip-%: %
	zip -j ./bin/waypoint-plugin-${PLUGIN_NAME}-$<_linux_amd64.zip ./bin/linux_amd64/waypoint-plugin-${PLUGIN_NAME}-$<
	zip -j ./bin/waypoint-plugin-${PLUGIN_NAME}-$<_darwin_amd64.zip ./bin/darwin_amd64/waypoint-plugin-${PLUGIN_NAME}-$<
	zip -j ./bin/waypoint-plugin-${PLUGIN_NAME}-$<_windows_amd64.zip ./bin/windows_amd64/waypoint-plugin-${PLUGIN_NAME}-$<.exe
	zip -j ./bin/waypoint-plugin-${PLUGIN_NAME}-$<_windows_386.zip ./bin/windows_386/waypoint-plugin-${PLUGIN_NAME}-$<.exe

# Build the plugin using a Docker container
build-docker:
	rm -rf ./releases
	DOCKER_BUILDKIT=1 docker build --build-arg VERSION=${VERSION} --output releases --progress=plain .
