package main

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/ChrisHeptagon/colibase/admin/db"
	"github.com/ChrisHeptagon/colibase/admin/server"
	"github.com/ChrisHeptagon/colibase/admin/utils"
)

func main() {
	if err := utils.HandleEnvVars(); err != nil {
		log.Fatalf("Error handling environment variables: %v", err)
	}
	dbChan := make(chan *sql.DB, 1)
	go func() {
		Db, err := db.InitDB()
		if err != nil {
			log.Fatalf("Error initializing DB: %v", err)
		}
		dbChan <- Db
	}()

	Db := <-dbChan
	defer Db.Close()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			err := Db.Ping()
			if err != nil {
				log.Printf("Database ping failed: %v", err)
				Db, err := db.InitDB()
				if err != nil {
					log.Fatalf("Error initializing DB: %v", err)
				}
				dbChan <- Db
			}
		}
	}()
	go func() {
		defer wg.Done()
		server.MainServer(Db)
	}()
	wg.Wait()
}
