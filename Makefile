PHONY: migrate_sqlite

migrate_sqlite:
	go run cmd/sqlite-migrator/main.go -storage-path=./storage/sso.db -migrations-path=./migrations -migrations-table=migrations

migrate_sqlite_down:
	go run cmd/sqlite-migrator/main.go -storage-path=./storage/sso.db -migrations-path=./migrations -migrations-table=migrations -down=true

migrate_sqlite_testdata:
	go run cmd/sqlite-migrator/main.go -storage-path=./storage/sso.db -migrations-path=./tests/migrations -migrations-table=testdatamigrations

migrate_sqlite_testdata_down:
	go run cmd/sqlite-migrator/main.go -storage-path=./storage/sso.db -migrations-path=./tests/migrations -migrations-table=testdatamigrations -down=true

run_test_app_server:
	go run cmd/sso/main.go --config=./config/test.json