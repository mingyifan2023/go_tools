package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

func main() {

	// 创建Gin实例  mojoru.com --- 给了一个4000端口
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(200, "hello https ")
	})

	// 设置路由
	localDir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前目录失败:", err)
		return
	}
	// autocert manager配置
	mng := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache(string(localDir)), // 缓存证书的文件夹

		HostPolicy: autocert.HostWhitelist("mojoru.com", "www.mojoru.com"),
	}

	// Gin服务器绑定自动证书管理器
	server := &http.Server{
		Addr:      ":https",
		Handler:   router,
		TLSConfig: &tls.Config{GetCertificate: mng.GetCertificate},
	}
	// Gin服务器绑定自动证书管理器
	log.Fatal(server.ListenAndServerTLS(":4000", ""))

	// // 设置监听的端口，以及SSL证书和私钥的位置
	// err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	// if err != nil {
	// 	log.Fatalf("启动HTTPS服务器失败: %v", err)
	// }
}
