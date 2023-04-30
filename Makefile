export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

os-archs=darwin:amd64 darwin:arm64 linux:amd64 linux:arm64

all:
	@$(foreach n, $(os-archs),\
		os=$(shell echo "$(n)" | cut -d : -f 1);\
		arch=$(shell echo "$(n)" | cut -d : -f 2);\
		gomips=$(shell echo "$(n)" | cut -d : -f 3);\
		target_suffix=$${os}_$${arch};\
		echo "Build $${os}-$${arch}...";\
		env CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} GOMIPS=$${gomips} go build -trimpath -ldflags "$(LDFLAGS)" -o ./release/codepass_$${target_suffix};\
		echo "Build $${os}-$${arch} done";\
	)

build:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o .

service: build
	./codepass service

clean:
	@rm -f ./codepass
	@rm -rf ./release
