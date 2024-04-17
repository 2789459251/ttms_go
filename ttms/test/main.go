package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/post", Upload)

	r.Run("0.0.0.0:8888")
}
