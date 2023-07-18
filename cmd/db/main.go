package main

import (
	"log"
	"net/url"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/clients"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/config"
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

func main() {
	c := config.GetConfig()

	db, err := clients.NewDBClient(c.DbUrl)
	if err != nil {
		panic(err)
	}

	url, err := obfuscatePassword(c.DbUrl)
	if err != nil {
		panic(err)
	}

	log.Printf("Applying schema to db: %s", url)
	err = db.ApplySchema()
	if err != nil {
		log.Printf("Error applying schema: %s", err.Error())
	}
}
