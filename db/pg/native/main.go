package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type User struct {
	ID          uuid.UUID
	UserName    string
	DisplayName string
}

func getUsersFromDB(db *sql.DB) ([]User, error) {
	var users []User
	rows, err := db.Query("SELECT id, username, display_name FROM public.users;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var idStr string
		var displayName sql.NullString
		var username sql.NullString

		err := rows.Scan(&idStr, &username, &displayName)
		if err != nil {
			return nil, err
		}
		id, _ := uuid.Parse(idStr)
		u := User{
			ID:          id,
			UserName:    username.String,
			DisplayName: displayName.String,
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
func main() {
	// 连接字符串，根据实际情况修改
	pgInfo := "postgres://root@localhost:26257/nakama?sslmode=disable"
	db, err := sql.Open("postgres", pgInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 确保数据库连接正常
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// 读取用户表到User列表
	users, err := getUsersFromDB(db)
	if err != nil {
		log.Fatal(err)
	}

	// 打印用户列表
	for _, user := range users {
		fmt.Printf("ID: %s, UserName: %s, DisplayName: %s\n", user.ID, user.UserName, user.DisplayName)
	}
}
