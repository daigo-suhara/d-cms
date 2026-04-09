package presenter

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Presenter handles rendering responses in either HTML or JSON format.
type Presenter interface {
	Render(c *gin.Context, status int, templateName string, data interface{})
	RenderError(c *gin.Context, status int, err error)
}

// Respond selects HTML or JSON presenter based on the Accept header.
func Respond(c *gin.Context) Presenter {
	accept := c.GetHeader("Accept")
	if strings.Contains(accept, "text/html") {
		return &HTMLPresenter{}
	}
	return &JSONPresenter{}
}
