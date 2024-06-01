package db

import (
	"flag"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var seed bool

func init() {
    reset := flag.Bool("r", false, "Reset DB")
    flag.BoolVar(&seed, "s", false, "Seed DB")

    flag.Parse()

    if *reset {
        os.Remove("./db/sqlite.db")
    }
}

func Init() *gorm.DB {
	dbConfig := &gorm.Config{
		TranslateError: true,
	}

	db, err := gorm.Open(sqlite.Open("./db/sqlite.db"), dbConfig)

	if err != nil {
		panic("Error connecting to the database")
	}

	db.AutoMigrate(
		&User{},
        &Event{},
	)

    if seed {
        Seed(db)
    }

    return db
}
