package main

import (
	"task_manager/server"
)

func main() {
	app := server.NewApp()
	app.Run()
}
