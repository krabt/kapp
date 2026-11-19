package main

import (
	"errors"
	"flag"
	"kapp/internal/data"
	"kapp/internal/data/sql"
	sqlDoc "kapp/internal/data/sql"
	"log"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

var (
	UserCommand string
	GenModel    bool
	GenTable    bool
	AutoMigrate bool
	GenTableDir = "sql"
	Args        []string
	dbPath      = "build/kapp.db"
)

func init() {
	flag.BoolVar(&GenModel, "g", false, "gorm gen model")
	flag.StringVar(&UserCommand, "u", "", "user command: list, create, reset")
	flag.StringVar(&dbPath, "d", dbPath, "SQLite database path")
	flag.BoolVar(&GenTable, "t", false, "goose cmd")
	flag.BoolVar(&AutoMigrate, "m", true, "apply database migrations before startup")
	flag.Parse()
	Args = flag.Args()

}

func sqliteDatabaseMissing(dsn string) bool {
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return false
	}
	if dsn == ":memory:" || strings.Contains(dsn, "mode=memory") {
		return false
	}
	path := strings.SplitN(dsn, "?", 2)[0]
	path = strings.TrimPrefix(path, "file:")
	path = strings.TrimPrefix(path, "//")
	_, err := os.Stat(path)
	return errors.Is(err, os.ErrNotExist)
}
func migrateAndSeed(dbPath string) {
	if err := data.GenTable(dbPath, sql.Files, []string{"up"}, ""); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

}

func main() {

	if GenModel {
		data.GenModel(dbPath)
		return
	}

	if GenTable {
		if len(Args) < 1 {
			flag.Usage()
			return
		}
		dir := ""
		if err := data.GenTable(dbPath, sqlDoc.Files, Args, dir); err != nil {
			log.Fatalf("failed to migrate database: %v", err)
		}

		return
	}
}
