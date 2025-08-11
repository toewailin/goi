package templates

// service_template.go - Template for generating service files
const ServiceTemplate = `package service

import (
	"{{.ModuleName}}/repository" // Import the repository layer dynamically
	"{{.ModuleName}}/models"     // Import your models (e.g., User)
	"fmt"
)

// {{.ServiceName}}Service provides business logic for {{.ServiceName}} operations
type {{.ServiceName}}Service struct {
	repo repository.{{.ServiceName}}Repository
}

// New{{.ServiceName}}Service creates a new instance of {{.ServiceName}}Service
func New{{.ServiceName}}Service(repo repository.{{.ServiceName}}Repository) *{{.ServiceName}}Service {
	return &{{.ServiceName}}Service{repo: repo}
}

// Create creates a new {{.ServiceName}} entity
func (s *{{.ServiceName}}Service) Create(name, email string) (*models.{{.ServiceName}}, error) {
	// Add business logic here if needed (e.g., validation)
	entity := &models.{{.ServiceName}}{Name: name, Email: email}
	return s.repo.Create(entity)
}

// GetByID retrieves a {{.ServiceName}} entity by its ID
func (s *{{.ServiceName}}Service) GetByID(id uint) (*models.{{.ServiceName}}, error) {
	return s.repo.FindByID(id)
}

// GetAll retrieves all {{.ServiceName}} entities
func (s *{{.ServiceName}}Service) GetAll() ([]*models.{{.ServiceName}}, error) {
	return s.repo.GetAll()
}

// Update updates an existing {{.ServiceName}} entity
func (s *{{.ServiceName}}Service) Update(id uint, name, email string) (*models.{{.ServiceName}}, error) {
	// Add business logic here if needed (e.g., validation)
	entity := &models.{{.ServiceName}}{ID: id, Name: name, Email: email}
	return s.repo.UpdateByID(id, entity)
}

// Delete deletes a {{.ServiceName}} entity by ID
func (s *{{.ServiceName}}Service) Delete(id uint) error {
	return s.repo.DeleteOne(id)
}
`
