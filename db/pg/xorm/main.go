package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

type User struct {
	ID          int64  `xorm:"'id' pk autoincr"` // 指定ID为自增主键
	UserName    string `xorm:"'user_name'"`
	DisplayName string `xorm:"'display_name'"`
}

func main() {
	// 连接字符串，根据实际情况修改
	dsn := "host=localhost port=26257 user=root dbname=xorm sslmode=disable"
	engine, err := xorm.NewEngine("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	println("New engine")

	// // 自动创建表
	// err = engine.Sync(new(User))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// println("Sync2")

	// // 初始化数据
	// for i := 0; i < 10; i++ {
	// 	u := User{UserName: "user" + strconv.Itoa(i), DisplayName: fmt.Sprintf("用户 %d", i)}
	// 	_, err = engine.Insert(&u)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// println("Init data")

	// 查询用户
	var users []User
	err = engine.Find(&users)
	if err != nil {
		log.Fatal(err)
	}
	println("Get all data")

	// 打印用户列表
	for _, user := range users {
		fmt.Printf("ID: %d, UserName: %s, DisplayName: %s\n", user.ID, user.UserName, user.DisplayName)
	}
	println("Print all data")
}
