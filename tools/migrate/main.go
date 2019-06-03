package main

import (
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

var DSN string
var Source string
var Step int

func init() {
	flag.StringVar(&DSN, "dsn", "", "dsn to a postgres database")
	flag.StringVar(&Source, "source", "", "migrations directory")
	flag.IntVar(&Step, "step", 0, "number of migrations to execute (0 is default and means all available)")
	flag.Parse()

	if DSN == "" {
		log.Println("dsn cannot be empty")
		os.Exit(1)
	}

	if Source == "" {
		log.Println("source cannot be empty")
		os.Exit(1)
	}

	fi, err := os.Stat(Source)
	if err != nil || !fi.IsDir() {
		log.Println("source is not a directory")
		os.Exit(1)
	}
}

func main() {
	m, err := migrate.New(Source, DSN)
	if err != nil {
		panic(err)
	}

	if err := m.Steps(Step); err != nil {
		if err == os.ErrNotExist {
			fmt.Println("no migrations to apply")
			os.Exit(0)
		}
		log.Println(err)
		os.Exit(1)
	}

	fmt.Println("migration complete")
}