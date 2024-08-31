package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Content string `json:"content"`
}

// 初始化数据库
func initDB() (*gorm.DB, error) {
	dbDir := "DB"
	// 创建 DB 文件夹（如果不存在）
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		os.Mkdir(dbDir, os.ModePerm)
	}

	dbPath := filepath.Join(dbDir, "todo.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移模式
	err = db.AutoMigrate(&Todo{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// 保存待办事项到数据库
func saveTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "不支持的方法", http.StatusMethodNotAllowed)
		return
	}

	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "无法解析请求体", http.StatusBadRequest)
		return
	}

	db, err := initDB()
	if err != nil {
		http.Error(w, "数据库错误", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer func() {
		db.Exec("PRAGMA foreign_keys = ON") // 确保启用外键支持
	}()

	// 插入数据
	result := db.Create(&todo)
	if result.Error != nil {
		http.Error(w, "插入数据失败", http.StatusInternalServerError)
		log.Fatal(result.Error)
		return
	}

	// 返回成功响应
	response := map[string]interface{}{
		"data": todo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 显示 HTML 模板
func indexHandler(w http.ResponseWriter, r *http.Request) {

	db, err := initDB()
	if err != nil {
		http.Error(w, "数据库错误", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer func() {
		db.Exec("PRAGMA foreign_keys = ON") // 确保启用外键支持
	}()
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, "无法加载模板", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func main() {
	// 初始化路由
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/save", saveTodoHandler)

	fmt.Println("服务器正在运行，访问地址: http://localhost:8004")
	log.Fatal(http.ListenAndServe(":8004", nil))
}
