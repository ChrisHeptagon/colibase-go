package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/ChrisHeptagon/colibase/admin/models"
	"github.com/ChrisHeptagon/colibase/admin/server"
	"github.com/ChrisHeptagon/colibase/admin/utils"
)

func main() {
	if err := utils.HandleEnvVars(); err != nil {
		log.Fatalf("Error handling environment variables: %v", err)
	}
	dbChan := make(chan *sql.DB, 1)
	go func() {
		Db, err := models.InitDB()
		if err != nil {
			log.Fatalf("Error initializing DB: %v", err)
		}
		dbChan <- Db
	}()

	Db := <-dbChan
	defer Db.Close()
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		wd, err := os.Getwd()
		fmt.Println(wd)
		if err != nil {
			log.Fatalf("Error getting working directory: %v", err)
		}
		cmd := exec.Command("node", fmt.Sprintf("%s/admin-ui/build/index.js", wd))
		cmd.Stdout = os.Stdout
		cmd.Env = os.Environ()
		err = cmd.Start()
		if err != nil {
			log.Fatalf("Error starting admin-ui: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			err := Db.Ping()
			if err != nil {
				log.Printf("Database ping failed: %v", err)
				Db, err := models.InitDB()
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
