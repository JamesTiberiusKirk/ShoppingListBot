package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/config"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/db"
)

func obfuscatePassword(connURL string) (string, error) {
	parsedURL, err := url.Parse(connURL)
	if err != nil {
		return "", err
	}

	username := parsedURL.User.Username()
	parsedURL.User = url.UserPassword(username, "xxxxxxx")

	return parsedURL.String(), nil
}

func applySchema(dbc *db.DB) {
	err := dbc.ApplySchema()
	if err != nil {
		log.Printf("Error applying schema: %s", err.Error())
		panic(err)
	}
}

func applyMigration(dbc *db.DB) {
	// type row struct {
	// 	ID      string `db:"id"`
	// 	Version int    `db:"version"`
	// }
	// var r []row
	// err := dbc.DB.Select(&r, "SELECT * FROM migrations")
	// if err != nil {
	// 	log.Printf("Error quering migrations table: %s", err.Error())
	// 	panic(err)
	// }
	//
	// lastVersion := r[0].Version
	//
}

func main() {
	c := config.GetConfig()

	dbc, err := db.NewDBClient(c.DbUrl)
	if err != nil {
		panic(err)
	}

	url, err := obfuscatePassword(c.DbUrl)
	if err != nil {
		panic(err)
	}

	action := flag.String("action", "bar", "provide actuion to be executed")
	flag.Parse()

	switch *action {
	case "schema":
		log.Printf("Applying schema to db: %s", url)
		applySchema(dbc)
	case "migration":
		log.Printf("Applying migration to db: %s", url)
		applyMigration(dbc)
	default:
		flag.PrintDefaults()
	}

}
