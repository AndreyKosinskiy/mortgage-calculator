run:
	cd cmd && go build -o ../build/backend.exe && cd ../build && backend.exe
build_image_old:
	docker build --tag mortgage-calculator --file ./deployments/Dockerfile .
run_container_old:
	docker run -d -p 8080:8080 mortgage-calculator
run_container:
	docker-compose -f ./deployments/docker-compose.yml up --build