package controller

import (
	"net/http"

	"github.com/daigo-suhara/d-cms/internal/presenter"
	"github.com/daigo-suhara/d-cms/internal/service"
	"github.com/gin-gonic/gin"
)

type APIKeyController struct {
	svc *service.APIKeyService
}

func NewAPIKeyController(svc *service.APIKeyService) *APIKeyController {
	return &APIKeyController{svc: svc}
}

func (ctrl *APIKeyController) Register(rg *gin.RouterGroup) {
	rg.GET("/api-keys", ctrl.List)
	rg.POST("/api-keys", ctrl.Create)
	rg.DELETE("/api-keys/:id", ctrl.Delete)
}

func (ctrl *APIKeyController) List(c *gin.Context) {
	keys, err := ctrl.svc.List(c.Request.Context())
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"status": 500, "error": err.Error()})
		return
	}
	data := gin.H{
		"title": "APIキー管理",
		"keys":  keys,
	}
	if flash, err := c.Cookie("flash_new_key"); err == nil {
		data["newKey"] = flash
		c.SetCookie("flash_new_key", "", -1, "/", "", false, true)
	}
	c.HTML(http.StatusOK, "api_keys/list.html", data)
}

func (ctrl *APIKeyController) Create(c *gin.Context) {
	name := c.PostForm("name")
	newKey, err := ctrl.svc.Create(c.Request.Context(), name)
	if err != nil {
		keys, _ := ctrl.svc.List(c.Request.Context())
		c.HTML(http.StatusUnprocessableEntity, "api_keys/list.html", gin.H{
			"title": "APIキー管理",
			"keys":  keys,
			"error": err.Error(),
		})
		return
	}
	c.SetCookie("flash_new_key", newKey.Key, 60, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/admin/api-keys")
}

func (ctrl *APIKeyController) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if c.GetHeader("HX-Request") == "true" {
		c.Status(http.StatusOK)
		return
	}
	presenter.Redirect(c, "/admin/api-keys")
}
