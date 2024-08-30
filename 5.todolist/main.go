package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Todo struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Content string `json:"content"`
}

type Database struct {
	Self *gorm.DB
}

// 初始化数据库
func (db *Database) Init() error {
	dbDir := "DB"
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err := os.Mkdir(dbDir, 0755)
		if err != nil {
			return fmt.Errorf("创建目录失败: %v", err)
		}
	}

	// 使用相对路径指定 SQLite 数据库文件位置
	dbFilePath := filepath.Join(dbDir, "todo.sqlite")

	// 初始化数据库
	_db, err := gorm.Open("sqlite3", dbFilePath)
	if err != nil {
		log.Printf("数据库连接失败. 数据库文件: %s", dbFilePath)
		return err
	}

	db.Self = _db

	// 自动迁移模式
	if err := db.Self.AutoMigrate(&Todo{}).Error; err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	return nil
}

func (db *Database) Close() {
	if err := db.Self.Close(); err != nil {
		log.Printf("关闭数据库时出错: %v", err)
	}
}

var DB Database

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

	// 插入数据
	result := DB.Self.Create(&todo)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("插入数据失败: %v", result.Error), http.StatusInternalServerError)
		log.Printf("插入数据失败: %v\n", result.Error)
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
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, "无法加载模板", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func main() {
	// 初始化数据库
	err := DB.Init()
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer DB.Close() // 在程序退出时关闭数据库

	// 初始化路由
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/save", saveTodoHandler)

	fmt.Println("服务器正在运行，访问地址: http://localhost:8004")
	log.Fatal(http.ListenAndServe(":8004", nil))
}
