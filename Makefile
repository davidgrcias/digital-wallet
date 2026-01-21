.PHONY: run stop test clean

run:
	@echo "Starting Database..."
	docker compose up -d
	@echo "Waiting for database to be ready..."
	@timeout /t 5 >nul 2>&1 || ping -n 5 127.0.0.1 >nul
	@echo "Starting Server..."
	go run cmd/api/main.go

stop:
	@echo "Stopping everything..."
	docker compose down

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning up..."
	docker compose down -v
	rm -f server
