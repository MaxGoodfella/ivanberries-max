build:
	docker-compose build --no-cache

run:
	docker-compose up -d

migrate:
	#docker-compose exec catalog-service goose -dir /app/catalog-service/db/migrations up
	docker-compose exec users-service goose -dir /app/users-service/db/migrations up

down:
	docker-compose down

restart:
	docker-compose down
	docker-compose up -d

#total restart: down build run migrate