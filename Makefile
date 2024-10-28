include .env

dev:
	@echo "Running the app..."
	@go run ./cmd/main.go
db:
	@echo "Initializing AD database..."
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres \
		-f ./migrations/init.sql
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d adsdb \
		-f ./migrations/create.sql
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d adsdb \
		-f ./migrations/insert.sql
	@echo "Has been successfully created"
build:
	@echo "Building the app, please wait..."
	@go build -o ./bin/ADs ./cmd/main.go
	@echo "Done."
build-cross:
	@echo "Bulding for windows, linux and macos, please wait..."
	@GOOS=linux GOARCH=amd64 go build -o ./bin/ADs-linux main.go
	@GOOS=darwin GOARCH=arm64 go build -o ./bin/ADs-macos main.go
	@GOOS=windows GOARCH=amd64 go build -o ./bin/ADs-windows main.go
	@echo "Done."
build-docker:
	@docker build -t ads-server .