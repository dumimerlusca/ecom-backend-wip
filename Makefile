include .env

.PHONY: run/api
run/api:
	go run ./cmd/api -port=${PORT} -db-dsn=${DB_DSN}


#############
# Migrations #
#############
.PHONY: migrations/new
migrations/new:
	migrate create -ext=.sql -dir ./migrations -seq $(name)

.PHONY: migrations/up
migrations/up:
	migrate -path ./migrations -database ${DB_DSN} -verbose up ${version}

.PHONY: migrations/down
migrations/down:
	migrate -path ./migrations -database ${DB_DSN} -verbose down ${version}