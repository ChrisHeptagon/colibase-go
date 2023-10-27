package main

import (
	"github.com/ChrisHeptagon/colibase/admin/db"
	"github.com/ChrisHeptagon/colibase/admin/server"
	"github.com/ChrisHeptagon/colibase/admin/utils"
)

func main() {
	utils.HandleEnvVars()
	Db, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	defer Db.Close()

	server.MainServer(Db)

}
