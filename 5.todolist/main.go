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
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		os.Mkdir(dbDir, os.ModePerm)
	}

	dbPath := filepath.Join(dbDir, "todo.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Todo{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// 查询待办事项
func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	db, err := initDB()
	if err != nil {
		http.Error(w, "数据库错误", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	var todos []Todo
	result := db.Find(&todos)
	if result.Error != nil {
		http.Error(w, "查询数据失败", http.StatusInternalServerError)
		log.Fatal(result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// 保存或更新待办事项到数据库
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

	// 检查 ID 是否存在以决定是插入新记录还是更新现有记录
	if todo.ID == 0 {
		result := db.Create(&todo)
		if result.Error != nil {
			http.Error(w, "插入数据失败", http.StatusInternalServerError)
			log.Fatal(result.Error)
			return
		}
	} else {
		result := db.Save(&todo) // 更新现有待办事项
		if result.Error != nil {
			http.Error(w, "更新数据失败", http.StatusInternalServerError)
			log.Fatal(result.Error)
			return
		}
	}

	response := map[string]interface{}{
		"data": todo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 显示 HTML 模板
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, "无法加载模板", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/save", saveTodoHandler)
	http.HandleFunc("/todos", getTodosHandler) // 获取待办事项的路由

	fmt.Println("服务器正在运行，访问地址: http://localhost:8004")
	log.Fatal(http.ListenAndServe(":8004", nil))
