package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName    string
	DisplayName string
}

func (u User) TableName() string {
	return "public.users"
}

func main() {
	// 连接字符串，根据实际情况修改
	dsn := "host=localhost port=26257 user=root dbname=gorm sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 自动迁移模式
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}

	// 初始化数据
	// for i := range 10 {
	// 	u := User{UserName: "user" + strconv.Itoa(i), DisplayName: fmt.Sprintf("用户 %d", i)}
	// 	db.Create(&u)
	// }

	// 查询用户
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	// 打印用户列表
	for _, user := range users {
		fmt.Printf("ID: %d, UserName: %s, DisplayName: %s, CreateAt: %v\n", user.ID, user.UserName, user.DisplayName, user.CreatedAt)
	}
}
