package main

import (
	"fmt"

	pb "github.com/defool/protoc-gen-gorm/example/generated/foo/v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	dsn := "root:123@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	checkErr(err)

	// Migrate the schema
	err = db.AutoMigrate(&pb.Group{}, &pb.Company{}, &pb.User{})
	checkErr(err)

	// Create
	u := &pb.User{Name: "def", Company: &pb.Company{Name: "cp1"}, Groups: []*pb.Group{{Name: "g1"}}}
	err = db.Create(u).Error
	checkErr(err)

	var user pb.User
	err = db.First(&user, 1).Error
	checkErr(err)
	fmt.Println("user", pb.DbFields.User.Name, "=", user.Name)
	fmt.Println("user", pb.DbFields.User.Id, "=", user.Id)
}
