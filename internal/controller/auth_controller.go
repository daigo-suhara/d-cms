package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	adminToken string
}

func NewAuthController(adminToken string) *AuthController {
	return &AuthController{adminToken: adminToken}
}

func (ctrl *AuthController) Register(rg *gin.RouterGroup) {
	rg.GET("/login", ctrl.LoginPage)
	rg.POST("/login", ctrl.Login)
	rg.POST("/logout", ctrl.Logout)
}

func (ctrl *AuthController) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/login.html", gin.H{})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	token := c.PostForm("token")
	if token != ctrl.adminToken {
		c.HTML(http.StatusUnauthorized, "auth/login.html", gin.H{
			"error": "Invalid token",
		})
		return
	}
	c.SetCookie("admin_token", token, 86400*7, "/", "", false, true)
	c.Redirect(http.StatusFound, "/admin/content-models")
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	c.SetCookie("admin_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/admin/login")
}
