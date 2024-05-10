migrate-local:
	@goose -dir db/migrations sqlite3 database.sqlite up