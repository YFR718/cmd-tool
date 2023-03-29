package api

import (
	"encoding/json"
	"github.com/YFR718/cmd-tool/server/cloud-disk/pkg/file-control"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func GetList(c *gin.Context) {
	path := c.Query("path")
	// 检查path是否为空
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "empty path",
		})
		return
	}
	// 创建文件对象
	f, err := file_control.NewFile(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	list, err := f.GetList()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	s, _ := json.Marshal(list)

	c.String(http.StatusOK, string(s))
}

func Mkdir(c *gin.Context) {
	path := c.Query("path")
	// 检查path是否为空
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "empty path",
		})
		return
	}
	err := file_control.MakeDir(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func RemoveFile(c *gin.Context) {
	path := c.Query("path")
	// 检查path是否为空
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "empty path",
		})
		return
	}
	// 创建文件对象
	f, err := file_control.NewFile(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = f.Remove()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func SendFile(c *gin.Context) {
	path := c.Query("path")

	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "empty path",
		})
		return
	}
	f, err := file_control.NewFile(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	if !f.Isdir {
		c.Header("Content-Disposition", "attachment; filename="+f.Name)
		c.Header("name", f.Name)
		c.File(f.Path)
		return
	}
	c.Header("zip", "true")
	c.Header("name", f.Name+".zip")
	c.Header("Content-Disposition", "attachment; filename="+f.Name+".zip")
	err = f.Zip()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "文件压缩失败" + err.Error(),
		})
		return
	}
	c.File(f.Path + ".zip")
	os.RemoveAll(f.Path + ".zip")
}

func GetFile(c *gin.Context) {
	path := c.Query("path")
	zip := c.Query("zip")

	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "empty path",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusInternalServerError, "读取file失败: "+err.Error())
		return
	}
	//fmt.Println("接收到文件: ", file.Filename, "存放路径：./data"+path+"/"+file.Filename)
	err = c.SaveUploadedFile(file, "./data"+path+"/"+file.Filename)
	if err != nil {
		c.String(http.StatusInternalServerError, "file保存失败: "+err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "success send " + file.Filename,
	})

	if zip == "true" {
		f, _ := file_control.NewFile(path + "/" + file.Filename)
		err = f.UnZip()
		if err != nil {
			c.String(http.StatusInternalServerError, "file解压失败: "+err.Error())
			return
		}
		os.RemoveAll("./data" + path + "/" + file.Filename)
	}

}
