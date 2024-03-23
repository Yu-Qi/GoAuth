package main

import (
	"fmt"

	"github.com/Yu-Qi/GoAuth/pkg/db"
	"github.com/Yu-Qi/GoAuth/pkg/db/model"
)

func tables() []interface{} {
	return []interface{}{
		&model.Account{},
	}
}

func create() error {
	// Create tables.
	if err := db.Get().AutoMigrate(tables()...); err != nil {
		return err
	}
	return nil
}

func main() {
	create()
	fmt.Println("Tables created.")
}
