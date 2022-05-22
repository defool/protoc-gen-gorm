package main

import (
	"fmt"

	pb "github.com/defool/protoc-gen-gorm/example/generated/foo/v1"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	checkErr(err)

	// Migrate the schema
	err = db.AutoMigrate(&pb.User{}, &pb.Group{}, &pb.Company{})
	checkErr(err)

	// Create
	err = db.Create(&pb.User{Name: "def"}).Error
	checkErr(err)

	var user pb.User
	err = db.First(&user, 1).Error
	checkErr(err)
	fmt.Println("user", pb.DbFields.User.Name, "=", user.Name)
}
