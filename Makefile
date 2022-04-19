ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run:
	docker-compose -f ./deployments/docker-compose.yml --env-file ./.env up --build
test:
	go test -v ./...
migrate:
	migrate -database postgres://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}/${DATABASE_NAME}?sslmode=${DATABASE_SSL_MODE} -path internal/migrations up
build_:
	cd cmd && go build -o ../build/backend.exe
run_:
	cd build && .\backend.exe
	