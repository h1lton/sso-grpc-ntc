package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(
		&storagePath,
		"storage-path",
		"",
		"путь к хранилищу",
	)
	flag.StringVar(
		&migrationsPath,
		"migrations-path",
		"",
		"путь к микрациям",
	)
	flag.StringVar(
		&migrationsTable,
		"migrations-table",
		"migrations",
		"имя таблицы миграций",
	)
	flag.Parse()

	if storagePath == "" {
		panic("Требуется путь к хранилищу")
	}
	if migrationsTable == "" {
		panic("Требуется путь к миграциям")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf(
			"sqlite3://%s?x-migrations-table=%s",
			storagePath,
			migrationsTable,
		),
	)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("нет миграций для применения")

			return
		}

		panic(err)
	}

	fmt.Println("миграции успешно применены")
}
