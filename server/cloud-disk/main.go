package main

import (
	"github.com/YFR718/cmd-tool/server/cloud-disk/internal/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.Default()
	r := routers.NewRouter()

	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
