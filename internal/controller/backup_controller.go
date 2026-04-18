package controller

import (
	"net/http"

	"github.com/daigo-suhara/d-cms/internal/service"
	"github.com/gin-gonic/gin"
)

type BackupController struct {
	svc *service.BackupService
}

func NewBackupController(svc *service.BackupService) *BackupController {
	return &BackupController{svc: svc}
}

func (ctrl *BackupController) Register(rg *gin.RouterGroup) {
	rg.GET("/backups", ctrl.List)
	rg.POST("/backups", ctrl.Create)
	rg.POST("/backups/restore", ctrl.Restore)
}

func (ctrl *BackupController) List(c *gin.Context) {
	backups, err := ctrl.svc.ListBackups(c.Request.Context())
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"status": 500, "error": err.Error()})
		return
	}
	data := gin.H{
		"title":   "バックアップ管理",
		"backups": backups,
	}
	if flash, err := c.Cookie("flash_backup"); err == nil {
		data["flash"] = flash
		c.SetCookie("flash_backup", "", -1, "/", "", false, true)
	}
	c.HTML(http.StatusOK, "backup/index.html", data)
}

func (ctrl *BackupController) Create(c *gin.Context) {
	key, err := ctrl.svc.CreateBackup(c.Request.Context())
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"status": 500, "error": err.Error()})
		return
	}
	c.SetCookie("flash_backup", "バックアップを作成しました: "+key, 60, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/admin/backups")
}

func (ctrl *BackupController) Restore(c *gin.Context) {
	key := c.PostForm("key")
	if key == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"status": 400, "error": "キーが指定されていません"})
		return
	}
	if err := ctrl.svc.RestoreBackup(c.Request.Context(), key); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"status": 500, "error": err.Error()})
		return
	}
	c.SetCookie("flash_backup", "バックアップからのインポートが完了しました", 60, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/admin/backups")
}
