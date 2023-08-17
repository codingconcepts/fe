postgres:
	docker run -d \
		--name postgres \
		-p 5432:5432 \
		-e POSTGRES_PASSWORD=password \
			postgres:15.2-alpine

postgres_create:
	PGPASSWORD=password psql -h localhost -p 5432 -d postgres -U postgres -f examples/postgres/create.sql

postgres_shell:
	PGPASSWORD=password psql -h localhost -p 5432 -d postgres -U postgres

run:
	go run fe.go postgres go \
		-u "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" \
		--go-package repo \
		-o examples/postgres/out.go

harness:
	go run examples/harness.go

clean:
	docker ps -aq | xargs docker rm -f