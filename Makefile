migrations-up:
	goose -dir ./migrations postgres "postgres://user:password@localhost:5432/news_db?sslmode=disable" up

migrations-down:
	goose -dir ./migrations postgres "postgres://user:password@localhost:5432/news_db?sslmode=disable" down