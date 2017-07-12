GOPATH ?= $(shell go env GOPATH)
GO_SRCS = $(shell find . -type f -name '*.go' -and -not -iwholename '*testdata*')
GO_PKGS = $(shell go list ./... | grep -v -e vendor -e fbs -e pb)
GO_VENDOR_PACKAGES = $(shell go list ./vendor/...)

GO_BUILD_FLAGS = -v -x
ifneq ($(GOD_RACE),)
	GO_BUILD_FLAGS+=-race
endif

build: bin/god

bin/god: $(GO_SRCS)
	go build $(GO_BUILD_FLAGS) -o $@ ./cmd/god

${GOPATH}/bin/god: $(GO_SRCS)
	go install $(GO_BUILD_FLAGS) ./cmd/god

install: ${GOPATH}/bin/god

clean:
	rm -rf ./bin

lint:
	@golint -set_exit_status $(GO_PKGS)

vet:
	@go vet -all -shadow $(GO_PKGS)

vendor/install:
	go install -v -x $(GO_VENDOR_PACKAGES)
	go install -v -x -race $(GO_VENDOR_PACKAGES)

pb: pb/fmt
	@rm -f $(shell find serial -type f -name '*.pb.go')
	protoc --gogofaster_out=plugins=grpc:. -I. -I$(GOPATH)/src $(shell find serial -type f -name '*.proto')

pb/fmt:
	clang-format -i $(shell find serial -type f -name '*.proto')


.PHONY: build install clean lint vet vendor/install pb pb/fmt
