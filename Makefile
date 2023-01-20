services-up:
	docker compose up -d postgres rabbitmq

services-down:
	docker compose down

migrate-create:
	docker run --user $(shell id -u):$(shell id -g) -v $(shell pwd)/db/migrations:/migrations --network host migrate/migrate \
	create -ext sql -dir /migrations -seq new-migration

migrate-up: services-up
	docker run --user $(shell id -u):$(shell id -g) -v $(shell pwd)/db/migrations:/migrations --network host migrate/migrate \
    -path=/migrations/ -database postgres://localhost:5651/database up

up:
	docker compose up

down:
	docker compose down -v

gen-docs:
	swagger generate spec -o ./docs/swagger.yml -i ./docs/swagger.yml -m