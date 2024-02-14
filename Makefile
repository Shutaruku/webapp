test:
	go test -v -cover -count=1 ./...

webappdb:
	docker run --name webappdb -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

createdb:
	docker exec -it webappdb createdb --username=postgres --owner=postgres webappdb

dropdb:
	docker exec -it webappdb dropdb --username=postgres webappdb

migrateup:
	soda migrate

migratedown:
	soda migrate down

.PHONY: test webappdb createdb dropdb migrateup migratedown