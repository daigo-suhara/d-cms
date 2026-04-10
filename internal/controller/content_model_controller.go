package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/daigo-suhara/d-cms/internal/domain"
	"github.com/daigo-suhara/d-cms/internal/presenter"
	"github.com/daigo-suhara/d-cms/internal/service"
	"github.com/gin-gonic/gin"
)

type ContentModelController struct {
	svc *service.ContentModelService
}

func NewContentModelController(svc *service.ContentModelService) *ContentModelController {
	return &ContentModelController{svc: svc}
}

func (ctrl *ContentModelController) Register(rg *gin.RouterGroup) {
	rg.GET("/content-models", ctrl.List)
	rg.GET("/content-models/new", ctrl.New)
	rg.POST("/content-models", ctrl.Create)
	rg.GET("/content-models/:id/edit", ctrl.Edit)
	rg.POST("/content-models/:id", ctrl.Update) // HTML forms don't support PUT
	rg.DELETE("/content-models/:id", ctrl.Delete)
}

func (ctrl *ContentModelController) List(c *gin.Context) {
	models, err := ctrl.svc.List(c.Request.Context())
	if err != nil {
		presenter.Respond(c).RenderError(c, http.StatusInternalServerError, err)
		return
	}
	presenter.Respond(c).Render(c, http.StatusOK, "content_models/list.html", gin.H{
		"models": models,
	})
}

func (ctrl *ContentModelController) New(c *gin.Context) {
	c.HTML(http.StatusOK, "content_models/form.html", gin.H{
		"model":  &domain.ContentModel{},
		"action": "create",
	})
}

func (ctrl *ContentModelController) Create(c *gin.Context) {
	m, err := ctrl.bindModel(c)
	if err != nil {
		presenter.Respond(c).RenderError(c, http.StatusBadRequest, err)
		return
	}

	if err := ctrl.svc.Create(c.Request.Context(), m); err != nil {
		status := httpStatus(err)
		if isHTMLRequest(c) {
			c.HTML(status, "content_models/form.html", gin.H{
				"model":  m,
				"action": "create",
				"error":  err.Error(),
			})
			return
		}
		presenter.Respond(c).RenderError(c, status, err)
		return
	}

	if isHTMLRequest(c) {
		presenter.Redirect(c, "/admin/content-models")
		return
	}
	c.JSON(http.StatusCreated, m)
}

func (ctrl *ContentModelController) Edit(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
		return
	}
	m, err := ctrl.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.HTML(httpStatus(err), "error.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "content_models/form.html", gin.H{
		"model":  m,
		"action": "update",
	})
}

func (ctrl *ContentModelController) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		presenter.Respond(c).RenderError(c, http.StatusBadRequest, err)
		return
	}
	m, err := ctrl.bindModel(c)
	if err != nil {
		presenter.Respond(c).RenderError(c, http.StatusBadRequest, err)
		return
	}
	m.ID = id

	if err := ctrl.svc.Update(c.Request.Context(), m); err != nil {
		status := httpStatus(err)
		if isHTMLRequest(c) {
			c.HTML(status, "content_models/form.html", gin.H{
				"model":  m,
				"action": "update",
				"error":  err.Error(),
			})
			return
		}
		presenter.Respond(c).RenderError(c, status, err)
		return
	}

	if isHTMLRequest(c) {
		presenter.Redirect(c, "/admin/content-models")
		return
	}
	c.JSON(http.StatusOK, m)
}

func (ctrl *ContentModelController) Delete(c *gin.Context) {
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

// bindModel parses ContentModel from either JSON body or multipart form.
func (ctrl *ContentModelController) bindModel(c *gin.Context) (*domain.ContentModel, error) {
	ct := c.ContentType()
	if ct == "application/json" {
		var m domain.ContentModel
		if err := c.ShouldBindJSON(&m); err != nil {
			return nil, fmt.Errorf("invalid JSON: %w", err)
		}
		return &m, nil
	}

	// HTML form submission
	name := c.PostForm("name")
	slug := c.PostForm("slug")
	fieldsJSON := c.PostForm("fields_json")

	var fields []domain.FieldDefinition
	if fieldsJSON != "" {
		if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
			return nil, fmt.Errorf("invalid fields JSON: %w", err)
		}
	}
	return &domain.ContentModel{
		Name:   name,
		Slug:   slug,
		Fields: fields,
	}, nil
}

// ── Public JSON API ──────────────────────────────────────────────────────────

func (ctrl *ContentModelController) RegisterAPI(rg *gin.RouterGroup) {
	rg.GET("/models", ctrl.APIList)
	rg.GET("/models/:slug", ctrl.APIGetBySlug)
}

func (ctrl *ContentModelController) APIList(c *gin.Context) {
	models, err := ctrl.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"models": models})
}

func (ctrl *ContentModelController) APIGetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	model, err := ctrl.svc.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(httpStatus(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, model)
}

func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err)
	}
	return uint(id), nil
}

func isHTMLRequest(c *gin.Context) bool {
	return c.GetHeader("HX-Request") == "true" || c.ContentType() == "application/x-www-form-urlencoded" || c.ContentType() == "multipart/form-data"
}
