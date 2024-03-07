package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"time"
)

func main() {
	// username and password are hardcoded because this should only be run on CI
	dbURL := "ci:github@tcp(localhost:3306)/migrate?parseTime=true"
	migrationsDir := "mysql/migrations"

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// CREATE DATABASE IF NOT EXISTS migrate;
	_, err = db.Exec(`CREATE DATABASE IF NOT EXISTS migrate;`)
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	m, err := migrate.New("file://"+migrationsDir, "mysql://"+dbURL)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			m.Down()
		}
	}()

	ver, err := validateMigrations(db, m)
	if err != nil {
		log.Fatalf("failed to validate migrations: %v", err)
		return
	}
	log.Printf("successfully validated %d migrations", ver)
	time.Sleep(1 * time.Second)
}

func validateMigrations(
	db *sql.DB,
	m *migrate.Migrate,
) (ver int, err error) {
	defer m.Down()

	var baseVer int
	var beforeDDL string
	for {
		beforeDDL, err = getDDLs(db)
		if err != nil {
			err = fmt.Errorf("failed to get beforeDDLs, ver: %d, err: %v", ver, err)
			return -1, err
		}
		err = m.Steps(1)
		if err != nil {
			break
		}
		version, _, _ := m.Version()
		log.Printf("base: %d,   up: %2d->%2d", baseVer, version-1, version)

		if err := m.Steps(-1); err != nil {
			log.Fatalf("failed to migrate down for version %d: %v", baseVer, err)
		}
		version, _, _ = m.Version()
		log.Printf("base: %d, down: %2d->%2d", baseVer, version+1, version)

		afterDDL, err := getDDLs(db)
		if err != nil {
			err = fmt.Errorf("failed to get afterDDLs, ver: %d, err: %v", baseVer, err)
			return baseVer, err
		}

		if beforeDDL != afterDDL {
			// show diff
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(beforeDDL, afterDDL, false)
			err = fmt.Errorf("failed to migrate down for version %d: DDLs are different, diffs: %v", baseVer, dmp.DiffPrettyText(diffs))
			return baseVer, err
		}

		if err := m.Steps(1); err != nil {
			err = fmt.Errorf("failed to migrate up after down for version %d: %v", baseVer, err)
			return baseVer, err
		}

		version, _, _ = m.Version()
		log.Printf("base: %d,   up: %2d->%2d", baseVer, version-1, version)

		baseVer++
	}

	if err.Error() != "file does not exist" {
		log.Printf("failed to proceed to next migration: ver: %d, err: %v", baseVer, err)
		return baseVer, err
	}
	return baseVer, nil
}

func getDDLs(db *sql.DB) (string, error) {
	tables, err := db.Query("SHOW TABLES")
	if err != nil {
		return "", err
	}
	defer tables.Close()

	var ddls string

	for tables.Next() {
		var tableName string
		if err := tables.Scan(&tableName); err != nil {
			return "", err
		}

		var createStatement string
		query := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
		err := db.QueryRow(query).Scan(&tableName, &createStatement)
		if err != nil {
			return "", err
		}
		ddls += fmt.Sprintf("Table: %s\nDDL:\n%s\n\n", tableName, createStatement)
	}

	return ddls, nil
}
