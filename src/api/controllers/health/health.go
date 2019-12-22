package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Alive(c *gin.Context)  {
	c.String(http.StatusOK, "Alive and Kicking")
}

