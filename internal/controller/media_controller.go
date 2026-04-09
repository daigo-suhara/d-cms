package controller

import (
	"net/http"

	"github.com/daigo-suhara/d-cms/internal/presenter"
	"github.com/daigo-suhara/d-cms/internal/service"
	"github.com/gin-gonic/gin"
)

type MediaController struct {
	svc *service.MediaService
}

func NewMediaController(svc *service.MediaService) *MediaController {
	return &MediaController{svc: svc}
}

func (ctrl *MediaController) Register(rg *gin.RouterGroup) {
	rg.GET("/media", ctrl.List)
	rg.POST("/media/upload", ctrl.Upload)
	rg.DELETE("/media/:id", ctrl.Delete)
}

func (ctrl *MediaController) List(c *gin.Context) {
	media, err := ctrl.svc.List(c.Request.Context())
	if err != nil {
		presenter.Respond(c).RenderError(c, http.StatusInternalServerError, err)
		return
	}
	presenter.Respond(c).Render(c, http.StatusOK, "media/list.html", gin.H{
		"media": media,
	})
}

func (ctrl *MediaController) Upload(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		presenter.Respond(c).RenderError(c, http.StatusBadRequest, err)
		return
	}
	m, err := ctrl.svc.Upload(c.Request.Context(), fh)
	if err != nil {
		presenter.Respond(c).RenderError(c, http.StatusInternalServerError, err)
		return
	}
	if isHTMLRequest(c) {
		presenter.Redirect(c, "/admin/media")
		return
	}
	c.JSON(http.StatusCreated, m)
}

func (ctrl *MediaController) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		presenter.Respond(c).RenderError(c, http.StatusBadRequest, err)
		return
	}
	if err := ctrl.svc.Delete(c.Request.Context(), id); err != nil {
		presenter.Respond(c).RenderError(c, httpStatus(err), err)
		return
	}
	if c.GetHeader("HX-Request") == "true" {
		c.Status(http.StatusOK)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
