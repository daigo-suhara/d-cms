package controller

import (
	"fmt"
	"net/http"

	"github.com/daigo-suhara/d-cms/internal/domain"
	"github.com/daigo-suhara/d-cms/internal/presenter"
	"github.com/daigo-suhara/d-cms/internal/service"
	"github.com/gin-gonic/gin"
)

type EntryController struct {
	svc      *service.EntryService
	modelSvc *service.ContentModelService
}

func NewEntryController(svc *service.EntryService, modelSvc *service.ContentModelService) *EntryController {
	return &EntryController{svc: svc, modelSvc: modelSvc}
}

func (ctrl *EntryController) RegisterAdmin(rg *gin.RouterGroup) {
	rg.GET("/:modelSlug/entries", ctrl.List)
	rg.GET("/:modelSlug/entries/new", ctrl.New)
	rg.POST("/:modelSlug/entries", ctrl.Create)
	rg.GET("/:modelSlug/entries/:id/edit", ctrl.Edit)
	rg.POST("/:modelSlug/entries/:id", ctrl.Update)
	rg.DELETE("/:modelSlug/entries/:id", ctrl.Delete)
}

func (ctrl *EntryController) RegisterAPI(rg *gin.RouterGroup) {
	rg.GET("/:modelSlug/entries", ctrl.APIList)
	rg.GET("/:modelSlug/entries/:id", ctrl.APIGet)
}

func (ctrl *EntryController) List(c *gin.Context) {
	modelSlug := c.Param("modelSlug")
	entries, model, err := ctrl.svc.ListBySlug(c.Request.Context(), modelSlug)
	if err != nil {
		presenter.Respond(c).RenderError(c, httpStatus(err), err)
		return
	}
	presenter.Respond(c).Render(c, http.StatusOK, "entries/list.html", gin.H{
		"model":   model,
		"entries": entries,
	})
}

func (ctrl *EntryController) New(c *gin.Context) {
	modelSlug := c.Param("modelSlug")
	model, err := ctrl.modelSvc.GetBySlug(c.Request.Context(), modelSlug)
	if err != nil {
		c.HTML(httpStatus(err), "error.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "entries/form.html", gin.H{
		"model":  model,
		"entry":  &domain.Entry{},
		"action": "create",
	})
}

func (ctrl *EntryController) Create(c *gin.Context) {
	modelSlug := c.Param("modelSlug")
	content, err := ctrl.parseContent(c, modelSlug)
	if err != nil {
		presenter.Respond(c).RenderError(c, http.StatusBadRequest, err)
		return
	}

	entry, err := ctrl.svc.Create(c.Request.Context(), modelSlug, content)
	if err != nil {
		status := httpStatus(err)
		if isHTMLRequest(c) {
			model, _ := ctrl.modelSvc.GetBySlug(c.Request.Context(), modelSlug)
			c.HTML(status, "entries/form.html", gin.H{
				"model":  model,
				"entry":  &domain.Entry{Content: content},
				"action": "create",
				"error":  err.Error(),
			})
			return
		}
		presenter.Respond(c).RenderError(c, status, err)
		return
	}

	if isHTMLRequest(c) {
		presenter.Redirect(c, fmt.Sprintf("/admin/%s/entries", modelSlug))
		return
	}
	c.JSON(http.StatusCreated, entry)
}

func (ctrl *EntryController) Edit(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
		return
	}
	entry, err := ctrl.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.HTML(httpStatus(err), "error.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "entries/form.html", gin.H{
		"model":  &entry.ContentModel,
		"entry":  entry,
		"action": "update",
	})
}

func (ctrl *EntryController) Update(c *gin.Context) {
	modelSlug := c.Param("modelSlug")
	id, err := parseID(c)
	if err != nil {
		presenter.Respond(c).RenderError(c, http.StatusBadRequest, err)
		return
	}
	content, err := ctrl.parseContent(c, modelSlug)
	if err != nil {
		presenter.Respond(c).RenderError(c, http.StatusBadRequest, err)
		return
	}

	entry, err := ctrl.svc.Update(c.Request.Context(), id, content)
	if err != nil {
		status := httpStatus(err)
		if isHTMLRequest(c) {
			model, _ := ctrl.modelSvc.GetBySlug(c.Request.Context(), modelSlug)
			c.HTML(status, "entries/form.html", gin.H{
				"model":  model,
				"entry":  &domain.Entry{ID: id, Content: content},
				"action": "update",
				"error":  err.Error(),
			})
			return
		}
		presenter.Respond(c).RenderError(c, status, err)
		return
	}

	if isHTMLRequest(c) {
		presenter.Redirect(c, fmt.Sprintf("/admin/%s/entries", modelSlug))
		return
	}
	c.JSON(http.StatusOK, entry)
}

func (ctrl *EntryController) Delete(c *gin.Context) {
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

// Public API handlers (no auth required)

func (ctrl *EntryController) APIList(c *gin.Context) {
	modelSlug := c.Param("modelSlug")
	entries, model, err := ctrl.svc.ListBySlug(c.Request.Context(), modelSlug)
	if err != nil {
		c.JSON(httpStatus(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"model":   model,
		"entries": entries,
	})
}

func (ctrl *EntryController) APIGet(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entry, err := ctrl.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(httpStatus(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entry)
}

// parseContent builds a ContentData map from either JSON body or form values.
func (ctrl *EntryController) parseContent(c *gin.Context, modelSlug string) (domain.ContentData, error) {
	if c.ContentType() == "application/json" {
		var body struct {
			Content domain.ContentData `json:"content"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			return nil, fmt.Errorf("invalid JSON: %w", err)
		}
		return body.Content, nil
	}

	// HTML form: values submitted as content[fieldName]
	model, err := ctrl.modelSvc.GetBySlug(c.Request.Context(), modelSlug)
	if err != nil {
		return nil, err
	}

	content := make(domain.ContentData)
	for _, field := range model.Fields {
		val := c.PostForm(fmt.Sprintf("content[%s]", field.Name))
		if val != "" {
			content[field.Name] = val
		}
	}
	return content, nil
}
