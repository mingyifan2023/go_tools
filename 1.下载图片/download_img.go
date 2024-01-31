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
		return fmt.Errorf("HTTP GET 请求失败:", err)
		
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("读取响应数据失败:", err)
		
	}
	// 如果文件已存在，先删除
	if _, err := os.Stat(filePath); err == nil {
		err = os.Remove(filePath)
		if err != nil {
			return fmt.Errorf("删除文件失败:", err)
			
		}
	}
	err = ioutil.WriteFile("/home/653aa6243a36b92bf4ce2fd2eab3e348.png", data, 0644)
	if err != nil {
		return  fmt.Errorf("保存文件失败:", err)
		
	}

	fmt.Println("图片已成功下载并保存")
	return nil 
}
