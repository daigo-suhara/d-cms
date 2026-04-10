package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"github.com/daigo-suhara/d-cms/config"
	"github.com/daigo-suhara/d-cms/internal/controller"
	"github.com/daigo-suhara/d-cms/internal/domain"
	"github.com/daigo-suhara/d-cms/internal/middleware"
	"github.com/daigo-suhara/d-cms/internal/repository"
	"github.com/daigo-suhara/d-cms/internal/service"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, storage service.StorageClient, cfg *config.Config) *gin.Engine {
	r := gin.Default()
	r.Static("/static", "./static")
	r.Use(middleware.RequestID())

	r.HTMLRender = createRenderer()

	// Repositories
	cmRepo := repository.NewContentModelRepository(db)
	entryRepo := repository.NewEntryRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)

	// Services
	cmSvc := service.NewContentModelService(cmRepo, entryRepo)
	entrySvc := service.NewEntryService(cmRepo, entryRepo)
	mediaSvc := service.NewMediaService(mediaRepo, storage)
	apiKeySvc := service.NewAPIKeyService(apiKeyRepo)

	// Controllers
	authCtrl := controller.NewAuthController(cfg.AdminToken)
	cmCtrl := controller.NewContentModelController(cmSvc)
	entryCtrl := controller.NewEntryController(entrySvc, cmSvc)
	mediaCtrl := controller.NewMediaController(mediaSvc)
	apiKeyCtrl := controller.NewAPIKeyController(apiKeySvc)

	// Root redirect
	r.GET("/", func(c *gin.Context) { c.Redirect(302, "/admin/content-models") })

	// Auth routes (public)
	adminPublic := r.Group("/admin")
	authCtrl.Register(adminPublic)

	// Protected admin routes
	admin := r.Group("/admin", middleware.Auth(cfg.AdminToken))
	cmCtrl.Register(admin)
	entryCtrl.RegisterAdmin(admin)
	mediaCtrl.Register(admin)
	apiKeyCtrl.Register(admin)

	// ── API v1 (全エンドポイントAPIキー必須) ────────────────────────────────

	api := r.Group("/api/v1", middleware.APIAuth(apiKeySvc))
	cmCtrl.RegisterAPI(api)
	entryCtrl.RegisterAPI(api)
	entryCtrl.RegisterAPIWrite(api)

	return r
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"fieldsJSON": func(fields []domain.FieldDefinition) template.JS {
			b, _ := json.Marshal(fields)
			return template.JS(b)
		},
		"jsonStringify": func(v any) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		},
		"dict": func(values ...interface{}) map[string]interface{} {
			m := make(map[string]interface{})
			for i := 0; i+1 < len(values); i += 2 {
				key, _ := values[i].(string)
				m[key] = values[i+1]
			}
			return m
		},
		"isImage": func(mimeType string) bool {
			return strings.HasPrefix(mimeType, "image/")
		},
		"formatSize": func(size int64) string {
			switch {
			case size >= 1<<20:
				return fmt.Sprintf("%.1f MB", float64(size)/(1<<20))
			case size >= 1<<10:
				return fmt.Sprintf("%.1f KB", float64(size)/(1<<10))
			default:
				return fmt.Sprintf("%d B", size)
			}
		},
		"keyPreview": func(key string) string {
			if len(key) <= 12 {
				return key
			}
			return key[:12] + "••••••••••••"
		},
	}
}

func createRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	fm := funcMap()

	layout := "templates/layout/base.html"
	partials := []string{
		"templates/partials/fields/text.html",
		"templates/partials/fields/number.html",
		"templates/partials/fields/date.html",
		"templates/partials/fields/markdown.html",
	}

	// Pages that use the main layout
	pages := map[string]string{
		"content_models/list.html": "templates/content_models/list.html",
		"content_models/form.html": "templates/content_models/form.html",
		"entries/list.html":        "templates/entries/list.html",
		"entries/form.html":        "templates/entries/form.html",
		"media/list.html":          "templates/media/list.html",
		"api_keys/list.html":       "templates/api_keys/list.html",
		"error.html":               "templates/error.html",
	}

	for name, page := range pages {
		files := append([]string{layout, page}, partials...)
		r.AddFromFilesFuncs(name, fm, files...)
	}

	// Auth pages (standalone, no layout)
	r.AddFromFilesFuncs("auth/login.html", fm, "templates/auth/login.html")

	return r
}
