// Usage: go run ./cmd/setadmin -db database.db -user <username>
package main

import (
	"flag"
	"fmt"
	"log"

	storm "github.com/asdine/storm/v3"
	"github.com/gtsteffaniak/filebrowser/backend/database/storage/bolt"
)

func main() {
	dbPath := flag.String("db", "database.db", "path to database.db")
	username := flag.String("user", "", "username to make admin")
	flag.Parse()

	if *username == "" {
		log.Fatal("usage: go run ./cmd/setadmin -user <username>")
	}

	db, err := storm.Open(*dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	store, err := bolt.NewStorage(db)
	if err != nil {
		log.Fatalf("new storage: %v", err)
	}

	// Always list users first
	all, err := store.Users.Gets()
	if err != nil {
		log.Fatalf("could not list users: %v", err)
	}
	fmt.Println("Users in database:")
	for _, u := range all {
		admin := ""
		if u.Permissions.Admin {
			admin = " [ADMIN]"
		}
		fmt.Printf("  %s%s\n", u.Username, admin)
	}

	if *username == "dummy" {
		return
	}

	user, err := store.Users.Get(*username)
	if err != nil {
		log.Fatalf("user %q not found: %v", *username, err)
	}

	user.Permissions.Admin = true
	if err := store.Users.Update(user, true, "Permissions"); err != nil {
		log.Fatalf("update failed: %v", err)
	}
	fmt.Printf("\n✓ %s is now admin\n", *username)
}
