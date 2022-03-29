PROJECT?=github.com/patrickjmcd/aws-iot-device-sdk-go
RELEASE?=$(shell git tag --sort=committerdate | tail -1)
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')


clean-cli:
	rm -f aws-provisioning-cli

clean-device:
	rm -f aws-provisioning-device

clean-device-arm7:
	rm -f aws-provisioning-device-arm7

clean-localproxy:
	rm -f localproxy

build-cli: clean-cli
	go build \
	-ldflags "-s -w -X github.com/patrickjmcd/go-version/version.Release=${RELEASE} \
	-X github.com/patrickjmcd/go-version/version.Commit=${COMMIT} -X github.com/patrickjmcd/go-version/version.BuildTime=${BUILD_TIME}" \
	-o aws-provisioning-cli \
	./cmd/cli


build-device-arm7: clean-device-arm7
	GOARCH=arm GOARM=7 GOOS=linux \
	go build \
  -ldflags "-s -w -X github.com/patrickjmcd/go-version/version.Release=${RELEASE} \
	-X github.com/patrickjmcd/go-version/version.Commit=${COMMIT} -X github.com/patrickjmcd/go-version/version.BuildTime=${BUILD_TIME}" \
	-o aws-provisioning-device-arm7 \
	./cmd/device 

build-device: clean-device
	go build \
  -ldflags "-s -w -X github.com/patrickjmcd/go-version/version.Release=${RELEASE} \
	-X github.com/patrickjmcd/go-version/version.Commit=${COMMIT} -X github.com/patrickjmcd/go-version/version.BuildTime=${BUILD_TIME}" \
	-o aws-provisioning-device \
	./cmd/device

build-localproxy: clean-localproxy
	go build \
	-ldflags "-s -w -X github.com/patrickjmcd/go-version/version.Release=${RELEASE} \
	-X github.com/patrickjmcd/go-version/version.Commit=${COMMIT} -X github.com/patrickjmcd/go-version/version.BuildTime=${BUILD_TIME}" \
	-o localproxy \
	./cmd/localproxy


clean-listen4tunnel:
	rm -f listen4tunnel*

ENDPOINT=$(shell cat ./.aws_endpoint.uri)
build-listen4tunnel: clean-listen4tunnel
ifeq ($(ENDPOINT),)
	@echo "*****"
	@echo "ENDPOINT not set"
	@echo "*****"
else
	go build \
	-ldflags "-s -w -X github.com/patrickjmcd/go-version/version.Release=${RELEASE} \
	-X github.com/patrickjmcd/go-version/version.Commit=${COMMIT} \
	-X github.com/patrickjmcd/go-version/version.BuildTime=${BUILD_TIME} \
	-X github.com/patrickjmcd/aws-iot-device-sdk-go/cmd/listen4tunnel/cfg.Endpoint=$(ENDPOINT)" \
	-o listen4tunnel \
	./cmd/listen4tunnel
endif

build-listen4tunnel-arm7: clean-listen4tunnel
ifeq ($(ENDPOINT),)
	@echo "*****"
	@echo "ENDPOINT not set"
	@echo "*****"
else
	GOARCH=arm GOARM=7 GOOS=linux \
	go build \
	-ldflags "-s -w -X github.com/patrickjmcd/go-version/version.Release=${RELEASE} \
	-X github.com/patrickjmcd/go-version/version.Commit=${COMMIT} \
	-X github.com/patrickjmcd/go-version/version.BuildTime=${BUILD_TIME} \
	-X github.com/patrickjmcd/aws-iot-device-sdk-go/cmd/listen4tunnel/cfg.Endpoint=$(ENDPOINT)" \
	-o listen4tunnel-arm7 \
	./cmd/listen4tunnel
endif

all: build-cli build-device build-device-arm7 build-localproxy build-listen4tunnel build-listen4tunnel-arm7