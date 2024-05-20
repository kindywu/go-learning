package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int32
	Name     string
	Password string
}

type UserProfile struct {
	UserID  int32
	Profile string
}

func insertUserAndProfile(db *sql.DB, user User, profile UserProfile) (userID int64, err error) {
	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// 插入User记录
	userInsertSQL := `INSERT INTO User (Name, Password) VALUES (?, ?)`
	result, err := tx.Exec(userInsertSQL, user.Name, user.Password)

	if err != nil {
		tx.Rollback() // 如果准备语句失败，则回滚事务
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected != 1 || err != nil {
		tx.Rollback() // 如果准备语句失败，则回滚事务
		return 0, err
	}

	// 获取自增ID
	userID, err = result.LastInsertId()
	if err != nil {
		tx.Rollback() // 如果获取ID失败，则回滚事务
		return 0, err
	}

	// 插入UserProfile记录
	profileInsertSQL := `INSERT INTO UserProfile (UserID, Profile) VALUES (?, ?)`
	result, err = tx.Exec(profileInsertSQL, userID, profile.Profile)
	if err != nil {
		tx.Rollback() // 如果执行插入失败，则回滚事务
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if rowsAffected != 1 || err != nil {
		tx.Rollback() // 如果准备语句失败，则回滚事务
		return 0, err
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return userID, nil
}

const createUserTableSQL = `CREATE TABLE IF NOT EXISTS User (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		password VARCHAR(32) NOT NULL
	)`
const createUserProfileTableSQL = `CREATE TABLE IF NOT EXISTS UserProfile (
		userID INT PRIMARY KEY,
		profile VARCHAR(255) NOT NULL
	)`

func main() {
	dsn := "root:@tcp(mysql:3306)/testdb2?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("DROP TABLE IF EXISTS User")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS UserProfile")
	if err != nil {
		log.Fatal(err)
	}

	// 迁移模式，创建表
	_, err = db.Exec(createUserTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(createUserProfileTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	user := User{Name: "John Doe", Password: "securepassword123"}
	profile := UserProfile{Profile: "John is a software developer."}

	userID, err := insertUserAndProfile(db, user, profile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User and profile inserted with user ID: %d\n", userID)
}
