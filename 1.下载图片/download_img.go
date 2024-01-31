package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "https://api.multiavatar.com/653aa6243a36b92bf4ce2fd2eab3e348.png"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("HTTP GET 请求失败:", err)
		return
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("读取响应数据失败:", err)
		return
	}

	err = ioutil.WriteFile("/home/653aa6243a36b92bf4ce2fd2eab3e348.png", data, 0644)
	if err != nil {
		fmt.Println("保存文件失败:", err)
		return
	}

	fmt.Println("图片已成功下载并保存")
}
