package main

import (
	"todo/db"
	"todo/handler"
)

func main() {
	err := db.DbTableInit()
	if err != nil {
		panic(err)
	}
	handler.InitServer()
}
