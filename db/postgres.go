package main

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/mcbenjemaa/go-stuff/internal/album"
)

func main() {
	example()
}

// Ref https://github.com/go-pg/pg/
// example of usage
func example() {
	// connect to db
	// pg.Connect(&pg.Options{
	//     User: "postgres",
	// })

	db := pg.Connect(&pg.Options{
		Addr:     ":5432",
		User:     "postgres",
		Password: "postgres",
		Database: "albums_db",
	})

	defer db.Close()

	// create schema
	err := createSchema(db)
	if err != nil {
		fmt.Println(fmt.Errorf("Error connecting to database, %w", err))
		panic(err)
	}
	fmt.Println("schema created!")

	// Insert Data
	alb1 := &album.Album{ID: "1", Title: "K8s", Artist: "Med Bj", Price: 100.4}
	_, err = db.Model(alb1).Insert()
	if err != nil {
		panic(err)
	}
	fmt.Println("album inserted!")

	// Select album by primary key.
	alb := &album.Album{ID: alb1.ID}
	err = db.Model(alb).WherePK().Select()
	if err != nil {
		panic(err)
	}

	// Select all albums.
	var albums []album.Album
	err = db.Model(&albums).Select()
	if err != nil {
		panic(err)
	}

	fmt.Printf("albums from db: %#v", albums)

}

// createSchema creates database schema for Album models.
func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*album.Album)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
