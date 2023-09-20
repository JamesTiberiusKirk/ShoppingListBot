package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/JamesTiberiusKirk/ShoppingListsBot/config"
	"github.com/JamesTiberiusKirk/ShoppingListsBot/db"
)

func sortArray(arr []int) []int {
	for i := 0; i <= len(arr)-1; i++ {
		for j := 0; j < len(arr)-1-i; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}

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
	type row struct {
		ID      string `db:"id"`
		Version int    `db:"version"`
	}
	var r row
	err := dbc.DB.QueryRowx("SELECT * FROM migrations WHERE id = 1").StructScan(&r)
	if err != nil {
		log.Printf("Error quering migrations table: %s", err.Error())
		panic(err)
	}
	log.Printf("Curent migration level: %d", r.Version)

	files, err := ioutil.ReadDir("./sql/migrations")
	if err != nil {
		log.Printf("Error opening migrations directory: %s", err.Error())
		panic(err)
	}

	var toApply []int

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		split := strings.Split(file.Name(), ".")
		level, err := strconv.Atoi(split[0])
		if err != nil {
			log.Printf("Could not parse migrations: %s", err.Error())
			panic(err)
		}

		if level > r.Version {
			toApply = append(toApply, level)
		}
	}

	if len(toApply) == 0 {
		log.Print("No new migrations")
		return
	}

	if len(toApply) > 1 {
		toApply = sortArray(toApply)
	}

	for _, l := range toApply {
		migration, err := os.ReadFile(fmt.Sprintf("./sql/migrations/%d.sql", l))
		if err != nil {
			log.Printf("Could not read migration file %d: %s", l, err.Error())
			panic(err)
		}

		tx, err := dbc.DB.Begin()
		if err != nil {
			panic(err)
		}

		_, err = tx.Exec(string(migration))
		if err != nil {
			panic(err)
		}

		_, err = tx.Exec(fmt.Sprintf(`
			INSERT INTO migrations (id, version)
			VALUES (1, %d)
			ON CONFLICT (id)
			DO UPDATE SET version = EXCLUDED.version;
		`, l))
		if err != nil {
			panic(err)
		}

		err = tx.Commit()
		if err != nil {
			log.Print("failed to commit transaction")
			panic(err)
		}

		log.Printf("Applied migration: %d", l)
		log.Printf("Upgraded migration version number: %d", l)

	}

}
func main() {
	log.Print("------------------------------------------------------------")
	log.Print("MIGRATOR")
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
	log.Print("------------------------------------------------------------")
}
