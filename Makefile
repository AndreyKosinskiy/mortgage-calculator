ifneq (,$(wildcard ./.env))
    include .env
    export
endif

migrate:
	migrate -database postgres://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}/${DATABASE_NAME}?sslmode=${DATABASE_SSL_MODE} -path internal/migrations up
run:
	docker-compose -f ./deployments/docker-compose.yml --env-file ./.env up --build
test:
	go test -v ./...
lint:
	golangci-lint run ./...
build_:
	cd cmd && go build -o ../build/backend.exe
run_:
	cd build && .\backend.exe
	