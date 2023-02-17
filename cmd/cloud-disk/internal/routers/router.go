package routers

import (
	"github.com/YFR718/cmd-tool/cmd/cloud-disk/internal/api"

	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	route := r.Group("/api")
	{
		// 获取文件夹下所有文件
		route.GET("/list", api.GetList)
		// 创建文件夹
		route.POST("/dir", api.Mkdir)
		// 上传文件
		route.POST("/file", api.GetFile)
		// 上传文件
		route.DELETE("/file", api.RemoveFile)
		// 下载文件
		route.GET("/file", api.SendFile)
	}
	return r

}
