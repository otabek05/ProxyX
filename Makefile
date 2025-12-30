APP_NAME=proxyx
ENTRY=./cmd/proxyx/

all: build

build:
	@mkdir -p dist/bin/linux
	GOOS=linux GOARCH=amd64 go build -o dist/bin/linux/$(APP_NAME) ./cmd/proxyx/

build-macos:
	@mkdir -p dist/bin/macos
	GOOS=darwin GOARCH=arm64 go build -o dist/bin/macos/$(APP_NAME) ./cmd/proxyx/
	@chmod +x bin/macos/$(APP_NAME)

build-all:
	GOOS=linux   GOARCH=amd64 go build -o dist/bin/linux/$(APP_NAME) $(ENTRY)
	GOOS=darwin  GOARCH=arm64 go build -o dist/bin/macos/$(APP_NAME) $(ENTRY)
	GOOS=windows GOARCH=amd64 go build -o dist/bin/windows/$(APP_NAME).exe $(ENTRY)


run:
	./bin/linux/$(APP_NAME)

install: build
	bash ./releaser/scripts/install_service.sh

uninstall:
	bash ./releaser/scripts/uninstall_service.sh

install-macos: build-macos
	bash ./releaser/scripts/install_service_macos.sh

uninstall-macos:
	bash ./releaser/scripts/uninstall_service_macos.sh

logs:
	sudo journalctl -u $(APP_NAME) -f
