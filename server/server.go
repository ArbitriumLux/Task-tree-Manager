package server

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// Server settings

const PostgresConnectionParameters = "host=localhost port=5432 user=postgres dbname=tree sslmode=disable password=root"

var Db *gorm.DB
var err error

//DataBaseConnection
func DataBaseConnection() {
	Db, err = gorm.Open("postgres", PostgresConnectionParameters)
	if err != nil {
		panic(err)
	}
}
