package presenter

import (
	"github.com/gin-gonic/gin"
)

type JSONPresenter struct{}

func (p *JSONPresenter) Render(c *gin.Context, status int, _ string, data interface{}) {
	c.JSON(status, data)
}

func (p *JSONPresenter) RenderError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}
