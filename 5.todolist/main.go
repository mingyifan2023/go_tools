package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
    "path/filepath"

    _ "github.com/mattn/go-sqlite3"
)

// Content 结构体用于解析请求中的 JSON 数据
type Content struct {
    Content string `json:"content"`
}

// 初始化数据库
func initDB() (*sql.DB, error) {
    dbDir := "./DB"
    // 创建 DB 文件夹（如果不存在）
    if _, err := os.Stat(dbDir); os.IsNotExist(err) {
        os.Mkdir(dbDir, os.ModePerm)
    }

    dbPath := filepath.Join(dbDir, "todo.db")
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    // 创建表格
    createTableSQL := `CREATE TABLE IF NOT EXISTS todos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT NOT NULL
    );`
    
    _, err = db.Exec(createTableSQL)
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

    var content Content
    err := json.NewDecoder(r.Body).Decode(&content)
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
    defer db.Close()

    // 插入数据
    stmt, err := db.Prepare("INSERT INTO todos(content) VALUES(?)")
    if err != nil {
        http.Error(w, "插入数据失败", http.StatusInternalServerError)
        log.Fatal(err)
        return
    }
    _, err = stmt.Exec(content.Content)
    if err != nil {
        http.Error(w, "插入数据失败", http.StatusInternalServerError)
        log.Fatal(err)
        return
    }

    // 返回成功响应
    response := map[string]interface{}{
        "data": content,
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
    // 初始化路由
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/save", saveTodoHandler)

    fmt.Println("服务器正在运行，访问地址: http://localhost:8004")
    log.Fatal(http.ListenAndServe(":8004", nil))
}
