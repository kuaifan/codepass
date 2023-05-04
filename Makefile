export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

os-archs=darwin:amd64 darwin:arm64 linux:386 linux:amd64 linux:arm linux:arm64 linux:mips64 linux:mips64le

all: assets
	@$(foreach n, $(os-archs),\
		os=$(shell echo "$(n)" | cut -d : -f 1);\
		arch=$(shell echo "$(n)" | cut -d : -f 2);\
		gomips=$(shell echo "$(n)" | cut -d : -f 3);\
		target_suffix=$${os}_$${arch};\
		echo "Build $${os}-$${arch}...";\
		env CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} GOMIPS=$${gomips} go build -trimpath -ldflags "$(LDFLAGS)" -o ./release/codepass_$${target_suffix};\
		echo "Build $${os}-$${arch} done";\
	)
	@cp ./release/codepass_linux_arm ./release/codepass_linux_aarch
	@cp ./release/codepass_linux_arm64 ./release/codepass_linux_aarch64

build: assets
	# go get github.com/jessevdk/go-assets-builder
    # go install github.com/jessevdk/go-assets-builder@latest
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o .

install: build
	./codepass install

service: build
	./codepass service

assets:
	go-assets-builder web/dist -o cmd/assets.go -p cmd

clean:
	@rm -f ./codepass
	@rm -rf ./release
