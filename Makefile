run:
	go mod tidy
	go run .
	
migrate.build:
	go build -o migrate migration/main/main.go

migrate.up:
	go run migration/main/main.go up

migrate.rollback:
	go run migration/main/main.go rollback