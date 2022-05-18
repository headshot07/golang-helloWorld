migrate-up:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/golang_project?sslmode=disable" -verbose up
migrate-down:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/golang_project?sslmode=disable" -verbose down
postgres-docker:
	docker run --name golang-postgres -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -e PGDATA=/var/lib/postgresql/data/pgdata -v ``$(PWD)``:/var/lib/postgresql/data -d postgres