package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type User struct {
	Username string
	Password string
}

func generateCookieStore() (*sessions.CookieStore, error) {
	store := sessions.NewCookieStore(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))
	return store, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求中的用户名和密码
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 这里应该有一个验证用户凭据的逻辑，例如查询数据库
	// 假设我们找到了一个匹配的用户
	validUser := User{Username: "admin", Password: "password"}

	if user.Username == validUser.Username && user.Password == validUser.Password {

		// 创建会话
		session, _ := sessionStore.Get(r, "user")

		data, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["user"] = string(data)
		err = sessionStore.Save(r, w, session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 返回成功响应
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
	} else {
		// 返回错误响应
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
	}
}

var sessionStore *sessions.CookieStore

func main() {
	// 在main函数中初始化sessionStore时添加错误检查
	var err error
	sessionStore, err = generateCookieStore()
	if err != nil {
		log.Fatal("Failed to generate session store:", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/login", loginHandler).Methods("POST")

	// 不需要登录的路由
	router.HandleFunc("/public", publicRoute).Methods("GET")

	// 需要登录的路由
	router.HandleFunc("/private", isAuthenticated(privateRoute)).Methods("GET")

	port := 9090
	fmt.Printf("Server started on :%d\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		log.Fatal(err)
	}
}

// publicRoute 处理不需要登录的请求
func publicRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "This is a public route. No login required."})
}

// ContextKey 是用于在上下文中存储认证信息的键
type ContextKey int

// 定义一个新的context键
const (
	CtxUserInfo ContextKey = iota
)

func isAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求中获取会话
		session, err := sessionStore.Get(r, "user")
		if err != nil {
			http.Error(w, "Error retrieving session", http.StatusInternalServerError)
			return
		}

		s, ok := session.Values["user"]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user := &User{}
		err = json.Unmarshal([]byte(s.(string)), user)
		if !ok {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), CtxUserInfo, user)
		r = r.WithContext(ctx)

		// 如果认证信息有效，继续处理请求
		next(w, r)
	}
}

// privateRoute 处理需要登录的请求
func privateRoute(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(CtxUserInfo).(User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// 继续处理需要登录的路由
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome back, " + user.Username + "! This is a private route."})
}
