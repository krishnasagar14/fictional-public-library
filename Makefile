.DEFAULT_GOAL := build
build: docker

docker:
	@docker-compose up --build -d