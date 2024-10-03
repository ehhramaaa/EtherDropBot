build:
	go build -o EtherDrop.exe

image:
	docker build -t etherdrop-bot .

up:
	docker-compose up -d

down:
	docker-compose down

delete:
	docker rmi etherdrop-bot --force

.PHONY: build up down delete