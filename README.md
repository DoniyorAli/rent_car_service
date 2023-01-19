migrate create -ext sql -dir ./storage/migrations -seq -digits 2 create_rentcar_table

migrate -path ./storage/migrations -database 'postgres://admin:qwerty123@localhost:5432/rentcar_service_db?sslmode=disable' up

migrate -path ./storage/migrations -database 'postgres://admin:qwerty123@localhost:5432/rentcar_service_db?sslmode=disable' down