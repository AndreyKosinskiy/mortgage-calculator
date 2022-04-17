run:
	docker-compose -f ./deployments/docker-compose.yml up --build
build_:
	cd cmd && go build -o ../build/backend.exe
run_:
	cd build && .\backend.exe
	