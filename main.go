package main

import (
	"Assignment/Packages/datastore"
	"Assignment/Packages/internet"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	datastore.RunMySQL()
}
func main() {
	defer datastore.CloseMySQL()
	internet.Htmlmain()
}
