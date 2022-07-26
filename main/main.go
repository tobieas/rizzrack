package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/initGit", initGit)
	r.GET("/initGit2", initGit2)
	r.GET("/cicd", cicd)

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
	r.Run(":9008")
}
