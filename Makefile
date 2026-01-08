.PHONY: all build run watch clean swag swag-init swag-fmt docker-up docker-down docker-logs docker-restart docker-clean

# Default: 
all: build

# Build
build:
	@go build -o ./tmp/main main.go

# Run 
run:
	@go run main.go

# Live Reload
dev:
	@AIR_PATH=$$(go env GOPATH)/bin/air; \
	if [ -f "$$AIR_PATH" ] || command -v air > /dev/null; then \
		if [ -f "$$AIR_PATH" ]; then \
			$$AIR_PATH; \
		else \
			air; \
		fi; \
		echo "Watching..."; \
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" = "y" ] || [ "$$choice" = "Y" ]; then \
			go install github.com/air-verse/air@latest; \
			$(go env GOPATH)/bin/air; \
			echo "Watching..."; \
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi

# Swagger documentation
swag:  swag-fmt swag-init

# Initialize Swagger documentation
swag-init:
	@swag init -g main.go -o ./docs

# Format Swagger documentation
swag-fmt:
	@swag fmt

# Clean
clean:
	@rm -rf ./tmp
	@echo "Cleaned build artifacts."

# Docker Compose commands
docker-up:
	@docker-compose up -d
	@echo "Docker services started!"

docker-down:
	@docker-compose down
	@echo "Docker services stopped!"

docker-logs:
	@docker-compose logs -f

docker-restart:
	@docker-compose restart
	@echo "Docker services restarted!"

docker-clean:
	@docker-compose down -v --remove-orphans
	@echo "Docker services and volumes cleaned!"