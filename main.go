// Command member-service owns library members: who they are and how to reach them.
package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	dsn := os.Getenv("DB_PATH")
	if dsn == "" {
		dsn = "members.db"
	}
	db, err := Open(dsn)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("member service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, NewRouter(db)); err != nil {
		log.Fatal(err)
	}
}
