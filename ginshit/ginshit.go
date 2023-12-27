package ginshit

import (
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		time.Sleep(6 * time.Second)
		c.String(200, "Welcome to Gin!")
	})

	// Use the pprof middleware
	pprof.Register(r)

	r.Run(":9999")
}
