package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ivanberries-max/internal/model"
	"ivanberries-max/internal/service/logic"
	"ivanberries-max/internal/service/validation/util"
	"net/http"
)

type CategoryHandler struct {
	service *logic.CategoryService
}

func NewCategoryHandler(service *logic.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(util.GetHTTPStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var category model.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong format"})
		return
	}

	if err := h.service.Create(&category); err != nil {
		c.JSON(util.GetHTTPStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong format"})
		return
	}

	category.ID = id

	updatedCategory, err := h.service.Update(&category)
	if err != nil {
		c.JSON(util.GetHTTPStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(util.GetHTTPStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "category deleted successfully"})
}
