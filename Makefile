migrateup:
	migrate -path db/migration -database "postgresql://root:4200@localhost:5432/film-db?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:4200@localhost:5432/film-db?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://root:4200@localhost:5432/film-db?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:4200@localhost:5432/film-db?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate
.PHONY: migrateup migrateup1 migratedown migratedown1 sqlc
