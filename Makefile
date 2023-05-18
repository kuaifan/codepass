export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w
GOEN := env CGO_ENABLED=0

os-archs=darwin:amd64 darwin:arm64 linux:386 linux:amd64 linux:arm linux:arm64 linux:mips64 linux:mips64le

all: assets
	@$(foreach n, $(os-archs),\
		os=$(shell echo "$(n)" | cut -d : -f 1);\
		arch=$(shell echo "$(n)" | cut -d : -f 2);\
		gomips=$(shell echo "$(n)" | cut -d : -f 3);\
		target_suffix=$${os}_$${arch};\
		echo "Build $${os}-$${arch}...";\
		$(GOEN) GOOS=$${os} GOARCH=$${arch} GOMIPS=$${gomips} go build -trimpath -ldflags "$(LDFLAGS)" -o ./release/codepass_$${target_suffix};\
		echo "Build $${os}-$${arch} done";\
	)
	@cp ./release/codepass_linux_arm ./release/codepass_linux_aarch
	@cp ./release/codepass_linux_arm64 ./release/codepass_linux_aarch64
	@cp ./release/codepass_linux_amd64 ./release/codepass_linux_x86_64

build: assets
	$(GOEN)  go build -trimpath -ldflags "$(LDFLAGS)" -o .

install: build
	./codepass install

.PHONY: service
service: build
	./codepass service --mode debug

.PHONY: assets
assets:
	$(GOEN) go-assets-builder shell -o assets/shell.go -p assets -v Shell
	$(GOEN) go-assets-builder web/dist -o assets/web.go -p assets -v Web

.PHONY: clean
clean:
	@rm -f ./codepass
	@rm -rf ./release


# 提示 go-assets-builder: No such file or directory 时解決辦法
# go get github.com/jessevdk/go-assets-builder
# go install github.com/jessevdk/go-assets-builder@latest