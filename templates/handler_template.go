package templates

// handler_template.go - Template for generating handlers
const HandlerTemplate = `package handlers

import (
	"fmt"
	"net/http"
	"{{.ModuleName}}/repository"  // Import the repository layer
	"{{.ModuleName}}/models"      // Import your models (e.g., User)
	"github.com/gin-gonic/gin"    // Import Gin framework
)

// {{.HandlerName}}Handler handles requests for {{.HandlerName}} resources
type {{.HandlerName}}Handler struct {
	repo repository.{{.HandlerName}}Repository
}

// New{{.HandlerName}}Handler creates a new instance of {{.HandlerName}}Handler
func New{{.HandlerName}}Handler(repo repository.{{.HandlerName}}Repository) *{{.HandlerName}}Handler {
	return &{{.HandlerName}}Handler{repo: repo}
}

// Index handles GET requests for {{.HandlerName}} resources
func (h *{{.HandlerName}}Handler) Index(c *gin.Context) {
	// Get all {{.HandlerName}} resources from the repository
	{{.HandlerName}}s, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve %s resources: %v", "{{.HandlerName}}", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": {{.HandlerName}}s})
}

// Show handles GET requests for a single {{.HandlerName}} resource by ID
func (h *{{.HandlerName}}Handler) Show(c *gin.Context) {
	// Get the ID from the URL parameters
	id := c.Param("id")

	// Fetch the resource by ID from the repository
	{{.HandlerName}}, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("%s resource not found: %v", "{{.HandlerName}}", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": {{.HandlerName}}})
}

// Create handles POST requests to create a new {{.HandlerName}} resource
func (h *{{.HandlerName}}Handler) Create(c *gin.Context) {
	var {{.HandlerName}} models.{{.HandlerName}}

	// Bind the request body to the {{.HandlerName}} struct
	if err := c.ShouldBindJSON(&{{.HandlerName}}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid input: %v", err)})
		return
	}

	// Call the repository to create the resource
	created{{.HandlerName}}, err := h.repo.Create(&{{.HandlerName}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create %s resource: %v", "{{.HandlerName}}", err)})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": created{{.HandlerName}}})
}

// Update handles PUT requests to update an existing {{.HandlerName}} resource
func (h *{{.HandlerName}}Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var {{.HandlerName}} models.{{.HandlerName}}

	// Bind the request body to the {{.HandlerName}} struct
	if err := c.ShouldBindJSON(&{{.HandlerName}}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid input: %v", err)})
		return
	}

	// Update the resource in the repository
	updated{{.HandlerName}}, err := h.repo.UpdateByID(id, &{{.HandlerName}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update %s resource: %v", "{{.HandlerName}}", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated{{.HandlerName}}})
}

// Delete handles DELETE requests to remove a {{.HandlerName}} resource by ID
func (h *{{.HandlerName}}Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Delete the resource from the repository
	if err := h.repo.DeleteOne(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete %s resource: %v", "{{.HandlerName}}", err)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
`
