export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

os-archs=darwin:amd64 darwin:arm64 linux:amd64 linux:arm64

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

build: assets
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o .

install: build
	./codepass install

service: build
	./codepass service --host=eeui.app --port=3443 --key=/Users/GAOYI/Downloads/eeui-app-nginx/app_key.key --crt=/Users/GAOYI/Downloads/eeui-app-nginx/app_chain.crt

assets:
	go-assets-builder web/dist -o cmd/assets.go -p cmd

clean:
	@rm -f ./codepass
	@rm -rf ./release
