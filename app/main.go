package app

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {

		//return response JSON
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	//mulai server dengan port 3000
	router.Run(":3000")

}
