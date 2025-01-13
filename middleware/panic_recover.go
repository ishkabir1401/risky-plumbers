package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func RecoverFromPanic(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			panicErrMsg := fmt.Errorf("%v", err)
			fmt.Println(c, "Panic detected", panicErrMsg, string(debug.Stack()))

			c.JSON(http.StatusInternalServerError, gin.H{"error": panicErrMsg})
		}
	}()

	c.Next()
}
