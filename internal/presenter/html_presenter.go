package presenter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTMLPresenter struct{}

func (p *HTMLPresenter) Render(c *gin.Context, status int, templateName string, data interface{}) {
	c.HTML(status, templateName, data)
}

func (p *HTMLPresenter) RenderError(c *gin.Context, status int, err error) {
	c.HTML(status, "error.html", gin.H{
		"status": status,
		"error":  err.Error(),
	})
}

// Redirect sends an HTMX-aware redirect: HX-Redirect for HTMX requests, standard 302 otherwise.
func Redirect(c *gin.Context, location string) {
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", location)
		c.Status(http.StatusNoContent)
		return
	}
	c.Redirect(http.StatusFound, location)
}
