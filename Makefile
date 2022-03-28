PROJECT?=github.com/patrickjmcd/aws-iot-devie-sdk-go


RELEASE?=$(shell git tag --sort=committerdate | tail -1)
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

DOCKER_REPO?=patrickjmcd
IMAGE_NAME?=in8n

clean-cli:
	rm -f aws-provisioning-cli

clean-device:
	rm -f aws-provisioning-device

clean-device-arm7:
	rm -f aws-provisioning-device-arm7

clean-localproxy:
	rm -f localproxy

clean-localproxy-arm7:
	rm -f localproxy-arm7

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

build-localproxy-arm7: clean-localproxy-arm7
	GOARCH=arm GOARM=7 GOOS=linux \
	go build \
	-ldflags "-s -w -X github.com/patrickjmcd/go-version/version.Release=${RELEASE} \
	-X github.com/patrickjmcd/go-version/version.Commit=${COMMIT} -X github.com/patrickjmcd/go-version/version.BuildTime=${BUILD_TIME}" \
	-o localproxy-arm7 \
	./cmd/localproxy


all: build-cli build-device build-device-arm7 build-localproxy build-localproxy-arm7