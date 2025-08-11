package templates

// repository_template.go - Template for generating repositories with CRUD operations
const RepositoryTemplate = `package repository

import (
	"fmt"
	"goi/models"  // Import your models (e.g., User)
	"gorm.io/gorm"
)

// {{.RepoName}}Repository defines the interface for CRUD operations
type {{.RepoName}}Repository interface {
	Create(item *models.{{.RepoName}}) (*models.{{.RepoName}}, error)
	CreateMany(items []*models.{{.RepoName}}) ([]*models.{{.RepoName}}, error)
	FindByID(id uint) (*models.{{.RepoName}}, error)
	GetAll() ([]*models.{{.RepoName}}, error)
	GetPaged(page, pageSize int) ([]*models.{{.RepoName}}, error)
	GetAllSorted(orderBy string, ascending bool) ([]*models.{{.RepoName}}, error)
	UpdateOne(item *models.{{.RepoName}}) (*models.{{.RepoName}}, error)
	UpdateByID(id uint, item *models.{{.RepoName}}) (*models.{{.RepoName}}, error)
	UpdateByEntity(item *models.{{.RepoName}}) (*models.{{.RepoName}}, error)
	UpdateMany(items []*models.{{.RepoName}}) ([]*models.{{.RepoName}}, error)
	SoftDelete(id uint) error
	DeleteOne(id uint) error
	DeleteMany(ids []uint) error
	Count() (int64, error)
}

// {{.RepoName}}Repo is the concrete implementation of {{.RepoName}}Repository
type {{.RepoName}}Repo struct {
	DB *gorm.DB
}

// New{{.RepoName}}Repository creates a new instance of {{.RepoName}}Repository
func New{{.RepoName}}Repository(db *gorm.DB) {{.RepoName}}Repository {
	return &{{.RepoName}}Repo{DB: db}
}

// Create adds a new {{.RepoName}} to the database
func (r *{{.RepoName}}Repo) Create(item *models.{{.RepoName}}) (*models.{{.RepoName}}, error) {
	if err := r.DB.Create(item).Error; err != nil {
		return nil, fmt.Errorf("failed to create {{.RepoName}}: %w", err)
	}
	return item, nil
}

// CreateMany adds multiple {{.RepoName}} entities to the database
func (r *{{.RepoName}}Repo) CreateMany(items []*models.{{.RepoName}}) ([]*models.{{.RepoName}}, error) {
	if err := r.DB.Create(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to create multiple {{.RepoName}}: %w", err)
	}
	return items, nil
}

// FindByID retrieves a {{.RepoName}} by its ID
func (r *{{.RepoName}}Repo) FindByID(id uint) (*models.{{.RepoName}}, error) {
	var item models.{{.RepoName}}
	if err := r.DB.First(&item, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find {{.RepoName}} by ID: %w", err)
	}
	return &item, nil
}

// GetAll retrieves all {{.RepoName}} entities from the database
func (r *{{.RepoName}}Repo) GetAll() ([]*models.{{.RepoName}}, error) {
	var items []*models.{{.RepoName}}
	if err := r.DB.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get all {{.RepoName}}: %w", err)
	}
	return items, nil
}

// GetPaged retrieves a paginated list of {{.RepoName}} entities
func (r *{{.RepoName}}Repo) GetPaged(page, pageSize int) ([]*models.{{.RepoName}}, error) {
	var items []*models.{{.RepoName}}
	offset := (page - 1) * pageSize
	if err := r.DB.Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get paginated {{.RepoName}}: %w", err)
	}
	return items, nil
}

// GetAllSorted retrieves all {{.RepoName}} entities sorted by a field
func (r *{{.RepoName}}Repo) GetAllSorted(orderBy string, ascending bool) ([]*models.{{.RepoName}}, error) {
	var items []*models.{{.RepoName}}
	if ascending {
		if err := r.DB.Order(orderBy).Find(&items).Error; err != nil {
			return nil, fmt.Errorf("failed to get sorted {{.RepoName}}: %w", err)
		}
	} else {
		if err := r.DB.Order(orderBy + " desc").Find(&items).Error; err != nil {
			return nil, fmt.Errorf("failed to get sorted {{.RepoName}}: %w", err)
		}
	}
	return items, nil
}

// UpdateOne updates a single {{.RepoName}} entity
func (r *{{.RepoName}}Repo) UpdateOne(item *models.{{.RepoName}}) (*models.{{.RepoName}}, error) {
	if err := r.DB.Save(item).Error; err != nil {
		return nil, fmt.Errorf("failed to update {{.RepoName}}: %w", err)
	}
	return item, nil
}

// UpdateByID updates a {{.RepoName}} entity by its ID
func (r *{{.RepoName}}Repo) UpdateByID(id uint, item *models.{{.RepoName}}) (*models.{{.RepoName}}, error) {
	if err := r.DB.Model(&models.{{.RepoName}}{}).Where("id = ?", id).Updates(item).Error; err != nil {
		return nil, fmt.Errorf("failed to update {{.RepoName}} by ID: %w", err)
	}
	return item, nil
}

// UpdateByEntity updates a {{.RepoName}} entity based on the entity itself
func (r *{{.RepoName}}Repo) UpdateByEntity(item *models.{{.RepoName}}) (*models.{{.RepoName}}, error) {
	if err := r.DB.Save(item).Error; err != nil {
		return nil, fmt.Errorf("failed to update {{.RepoName}} by entity: %w", err)
	}
	return item, nil
}

// UpdateMany updates multiple {{.RepoName}} entities
func (r *{{.RepoName}}Repo) UpdateMany(items []*models.{{.RepoName}}) ([]*models.{{.RepoName}}, error) {
	if err := r.DB.Save(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to update multiple {{.RepoName}}: %w", err)
	}
	return items, nil
}

// SoftDelete marks a {{.RepoName}} entity as deleted (without actually deleting it)
func (r *{{.RepoName}}Repo) SoftDelete(id uint) error {
	if err := r.DB.Model(&models.{{.RepoName}}{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
		return fmt.Errorf("failed to soft delete {{.RepoName}}: %w", err)
	}
	return nil
}

// DeleteOne deletes a single {{.RepoName}} by its ID
func (r *{{.RepoName}}Repo) DeleteOne(id uint) error {
	if err := r.DB.Delete(&models.{{.RepoName}}{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete {{.RepoName}} by ID: %w", err)
	}
	return nil
}

// DeleteMany deletes multiple {{.RepoName}} entities by their IDs
func (r *{{.RepoName}}Repo) DeleteMany(ids []uint) error {
	if err := r.DB.Where("id IN ?", ids).Delete(&models.{{.RepoName}}{}).Error; err != nil {
		return fmt.Errorf("failed to delete multiple {{.RepoName}} by IDs: %w", err)
	}
	return nil
}

// Count returns the total number of {{.RepoName}} records in the database
func (r *{{.RepoName}}Repo) Count() (int64, error) {
	var count int64
	if err := r.DB.Model(&models.{{.RepoName}}{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count {{.RepoName}}: %w", err)
	}
	return count, nil
}
`
