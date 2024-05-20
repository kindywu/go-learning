package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// User 模型对应数据库中的User表
type User struct {
	ID        int64            `json:"id"`
	Name      string           `json:"name"`
	Password  string           `json:"password"`
	Avatar    sql.Null[string] `json:"avatar"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	Version   int              `json:"version"`
}

const createTableSQL = `CREATE TABLE IF NOT EXISTS user (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		password VARCHAR(32) NOT NULL,
		avatar VARCHAR(255), 
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		version INT DEFAULT 1
	)`

// createUser SQL 插入语句
const createUserSQL = `INSERT INTO user (name, password, avatar, created_at, updated_at, version) VALUES (?, ?, ?, NOW(), NOW(), 1)`

// getUserByIDSQL SQL 查询语句
const getUserByIDSQL = `SELECT id, name, password, avatar, created_at, updated_at, version FROM user WHERE id = ?`

// updateUserSQL SQL 更新语句
const updateUserSQL = `UPDATE user SET name = ?, password = ?, avatar = ?, updated_at = NOW(), version = version + 1 WHERE id = ? AND version = ?`

// deleteUserSQL SQL 删除语句
const deleteUserSQL = `DELETE FROM user WHERE id = ?`

func main() {
	// 连接到MySQL数据库
	dsn := "root:@tcp(mysql:3306)/testdb?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 迁移模式，创建表
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// 示例CRUD操作
	// 创建用户
	user := User{Name: "admin", Password: "admin123"}
	hashedPassword := md5AndHex(user.Password)
	user.Avatar = sql.Null[string]{}
	_, err = db.Exec(createUserSQL, user.Name, hashedPassword, user.Avatar)
	if err != nil {
		log.Fatal(err)
	}

	// 获取用户
	id := int64(1) // 假设我们要查询的用户ID为1
	var fetchedUser User
	// var avatar sql.NullString
	err = db.QueryRow(getUserByIDSQL, id).Scan(&fetchedUser.ID, &fetchedUser.Name, &fetchedUser.Password, &fetchedUser.Avatar, &fetchedUser.CreatedAt, &fetchedUser.UpdatedAt, &fetchedUser.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return // 返回nil表示没有找到记录
		}
		log.Fatal(err)
	}

	// if avatar.Valid {
	// 	fetchedUser.Avatar = &avatar.String
	// }
	fmt.Printf("Fetched User: %+v\n", fetchedUser)

	// 更新用户
	newPassword := "newpassword"
	updatedUser, err := db.Exec(updateUserSQL, user.Name, md5AndHex(newPassword), user.Avatar, user.ID, user.Version)
	if err != nil {
		log.Fatal(err)
	}
	if affected, err := updatedUser.RowsAffected(); err != nil {
		fmt.Printf("Rows affected: %d\n", affected)
	}

	// 删除用户
	deletedUser, err := db.Exec(deleteUserSQL, id)
	if err != nil {
		log.Fatal(err)
	}
	if affected, err := deletedUser.RowsAffected(); err != nil {
		fmt.Printf("Rows affected: %d\n", affected)
	}
}

// md5AndHex 用于生成MD5加密后的字符串并转换为十六进制
func md5AndHex(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}
