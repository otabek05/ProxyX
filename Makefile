APP_NAME=proxyx

all: build

build:
	@mkdir -p bin/linux
	GOOS=linux GOARCH=amd64 go build -o bin/linux/$(APP_NAME) ./cmd/proxy/
	@chmod +x bin/linux/$(APP_NAME)

build-macos:
	@mkdir -p bin/macos
	GOOS=darwin GOARCH=arm64 go build -o bin/macos/$(APP_NAME) ./cmd/proxy/
	@chmod +x bin/macos/$(APP_NAME)

run:
	./bin/linux/$(APP_NAME)

install: build
	bash ./scripts/install_service.sh

uninstall:
	bash ./scripts/uninstall_service.sh

install-macos: build-macos
	bash ./scripts/install_service_macos.sh

uninstall-macos:
	bash ./scripts/uninstall_service_macos.sh

logs:
	sudo journalctl -u $(APP_NAME) -f
