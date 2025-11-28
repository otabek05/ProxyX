APP_NAME=proxyx
BIN=bin/$(APP_NAME)

all: build

build:
	@echo "Building $(APP_NAME)..."
	GOOS=linux GOARCH=amd64 go build -o $(BIN) ./cmd/proxy/
	@chmod +x $(BIN)
	@echo "Build complete → $(BIN) is now executable"

run: build
	@echo "Running $(APP_NAME)..."
	./bin/proxy --config=configs/proxy.yaml --port=8000

install: build
	@echo "Installing service..."
	bash ./scripts/install_service.sh

uninstall:
	@echo "Uninstalling service..."
	bash ./scripts/uninstall_service.sh

start:
	sudo systemctl start $(APP_NAME)

stop:
	sudo systemctl stop $(APP_NAME)

restart:
	sudo systemctl restart $(APP_NAME)

status:
	sudo systemctl status $(APP_NAME)

logs:
	sudo journalctl -u $(APP_NAME) -f
