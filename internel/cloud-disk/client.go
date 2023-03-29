package cloud_disk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	file_control "github.com/YFR718/cmd-tool/server/cloud-disk/pkg/file-control"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type File struct {
	Path  string
	Name  string
	Isdir bool
	data  []byte
}

var client = &http.Client{}

func getlist(path string, print bool) error {
	// 创建一个 HTTP 请求
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/api/list?path="+path, nil)
	if err != nil {
		return err
	}

	// 发送 HTTP 请求
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	s := make([]byte, 1000)
	n, _ := res.Body.Read(s)
	if res.Status != "200 OK" {
		fmt.Println(res.Status)
		return errors.New(string(s[:n]))
	}

	if print {
		var files []File
		json.Unmarshal(s[:n], &files)
		for _, file := range files {
			if file.Isdir {
				fmt.Printf("\u001B[34m %v \u001B[0m ", file.Name)
			} else {
				fmt.Printf("\u001B[32m %v \u001B[0m ", file.Name)
			}
		}
		fmt.Println("")
	}

	return nil
}

//fmt.Println("\033[31mThis text is red.\033[0m")
//fmt.Println("\033[32mThis text is green.\033[0m")
//fmt.Println("\033[33mThis text is yellow.\033[0m")
//fmt.Println("\033[34mThis text is blue.\033[0m")
//fmt.Println("\033[35mThis text is magenta.\033[0m")
//fmt.Println("\033[36mThis text is cyan.\033[0m")

type ProgressPrinter struct {
	total    int64
	progress int64
}

func (p *ProgressPrinter) Write(b []byte) (n int, err error) {
	n = len(b)
	p.progress += int64(n)
	fmt.Printf("\rUploading... %d%%", p.progress*100/p.total)
	if p.progress == p.total {
		fmt.Println()
	}
	return
}

type ProgressPrinter2 struct {
	total int
}

func (p *ProgressPrinter2) Write(b []byte) (n int, err error) {
	n = len(b)
	p.total++
	switch p.total % 5 {
	case 0:
		fmt.Printf("\rDownLoading.")
	case 1:
		fmt.Printf("\rDownLoading..")
	case 2:
		fmt.Printf("\rDownLoading...")
	case 3:
		fmt.Printf("\rDownLoading....")
	case 4:
		fmt.Printf("\rDownLoading.....")
	}

	return
}

func pushFile(path, local string) error {
	fmt.Println("path and local:", path, local)

	filestate, err := os.Stat(local)
	if err != nil {
		return err
	}
	zip := false
	if filestate.IsDir() {
		//path = filepath.Join(path, filestate.Name())
		//
		//mkdir(path)
		zip = true
		f := file_control.File{Path: local}
		println("p1")
		err = f.Zip()
		println("p2", err)
		if err != nil {
			return err
		}
		local += ".zip"
		defer os.RemoveAll(local)
	}

	//// 创建一个 multipart 请求体
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)
	//// 打开文件
	file, err := os.Open(local)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// 创建一个 multipart 表单域，将文件内容写入其中
	fileField, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return err
	}
	_, err = io.Copy(fileField, file)
	if err != nil {
		return err
	}

	// 关闭 multipart 请求体
	writer.Close()

	url := "http://127.0.0.1:8080/api/file?path=" + path + "&zip=false"
	if zip {
		url = "http://127.0.0.1:8080/api/file?path=" + path + "&zip=true"
	}
	// 创建一个 HTTP 请求 requestBody
	request, err := http.NewRequest("POST", url, io.TeeReader(requestBody, &ProgressPrinter{total: stat.Size()}))
	if err != nil {
		return err
	}
	// 设置请求头
	request.Header.Set("Content-Type", writer.FormDataContentType())
	//request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundaryXXXXXX")

	// 发送 HTTP 请求并等待响应
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// 读取响应体
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// 打印响应体
	fmt.Println("\n", string(responseBytes))
	return nil
}

func pullFile(path string) error {
	// 创建一个 HTTP 请求
	fmt.Println("getting file... ")
	request, err := http.NewRequest("GET", "http://127.0.0.1:8080/api/file?path="+path, nil)
	if err != nil {
		return err
	}

	// 发送 HTTP 请求并等待响应
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	zip := response.Header.Get("zip")
	name := response.Header.Get("name")

	// 创建一个本地文件
	twd, _ := os.Getwd()
	local := filepath.Join(twd, name)
	file, err := os.Create(local)
	if err != nil {
		return err
	}

	// 将服务器返回的文件写入本地文件,并动态打印进度条
	_, err = io.Copy(file, io.TeeReader(response.Body, &ProgressPrinter2{total: 0}))
	if err != nil {
		return err
	}

	//_, err = io.Copy(file, response.Body)

	fmt.Println("\n File downloaded successfully: ", name)

	// 解压缩
	if zip == "true" {
		f := file_control.File{Path: local}
		if err != nil {
			println(err.Error())
		}

		err = f.UnZip()
		if err != nil {
			return err
		}

		file.Close()
		os.RemoveAll(local)
	}
	return nil
}

func mkdir(path string) error {
	// 创建一个 HTTP 请求
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/api/dir?path="+path, nil)
	if err != nil {
		return err
	}

	// 发送 HTTP 请求
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	s := make([]byte, 1000)
	n, _ := res.Body.Read(s)
	if res.Status != "200 OK" {
		fmt.Println(res.Status)
		return errors.New(string(s[:n]))
	}

	return nil
}

func rmFile(path string) error {
	// 创建一个 HTTP 请求
	req, err := http.NewRequest("DELETE", "http://127.0.0.1:8080/api/file?path="+path, nil)
	if err != nil {
		return err
	}

	// 发送 HTTP 请求
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	s := make([]byte, 1000)
	n, _ := res.Body.Read(s)
	if res.Status != "200 OK" {
		fmt.Println(res.Status)
		return errors.New(string(s[:n]))
	}

	return nil
}
