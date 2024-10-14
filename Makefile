build:
	@go build -o bin/app cmd/main.go

run: build
	@./bin/app

migration:
	@migrate create -ext sql -dir internal/db/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run  /app/cmd/migrate/main.go up

migrate-down:
	@go run app/cmd/migrate/main.go down


