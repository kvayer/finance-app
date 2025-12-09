.PHONY: run build stop

run:
	docker-compose up --build

stop:
	docker-compose down

build:
	go build -o ./.bin/app cmd/app/main.go