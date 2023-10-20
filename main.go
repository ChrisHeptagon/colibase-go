package main

import (
	"github.com/ChrisHeptagon/colibase/db"
	"github.com/ChrisHeptagon/colibase/server"
	"github.com/ChrisHeptagon/colibase/utils"
)

func main() {
	utils.HandleEnvVars()
	db, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	server.MainServer(db)
}
