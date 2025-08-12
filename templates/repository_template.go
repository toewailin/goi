package templates

// repository_template.go - Template for generating repositories with CRUD operations
const RepositoryTemplate = `package repository

import (
	"fmt"
	"goi/models"  // Import your models (e.g., User)
	"gorm.io/gorm"
)

// {{.RepositoryName}}Repository defines the interface for CRUD operations
type {{.RepositoryName}}Repository interface {
	Create(item *models.{{.RepositoryName}}) (*models.{{.RepositoryName}}, error)
	CreateMany(items []*models.{{.RepositoryName}}) ([]*models.{{.RepositoryName}}, error)
	FindByID(id uint) (*models.{{.RepositoryName}}, error)
	GetAll() ([]*models.{{.RepositoryName}}, error)
	GetPaged(page, pageSize int) ([]*models.{{.RepositoryName}}, error)
	GetAllSorted(orderBy string, ascending bool) ([]*models.{{.RepositoryName}}, error)
	UpdateOne(item *models.{{.RepositoryName}}) (*models.{{.RepositoryName}}, error)
	UpdateByID(id uint, item *models.{{.RepositoryName}}) (*models.{{.RepositoryName}}, error)
	UpdateByEntity(item *models.{{.RepositoryName}}) (*models.{{.RepositoryName}}, error)
	UpdateMany(items []*models.{{.RepositoryName}}) ([]*models.{{.RepositoryName}}, error)
	SoftDelete(id uint) error
	DeleteOne(id uint) error
	DeleteMany(ids []uint) error
	Count() (int64, error)
}

// {{.RepositoryName}}Repo is the concrete implementation of {{.RepositoryName}}Repository
type {{.RepositoryName}}Repo struct {
	DB *gorm.DB
}

// New{{.RepositoryName}}Repository creates a new instance of {{.RepositoryName}}Repository
func New{{.RepositoryName}}Repository(db *gorm.DB) {{.RepositoryName}}Repository {
	return &{{.RepositoryName}}Repo{DB: db}
}

// Create adds a new {{.RepositoryName}} to the database
func (r *{{.RepositoryName}}Repo) Create(item *models.{{.RepositoryName}}) (*models.{{.RepositoryName}}, error) {
	if err := r.DB.Create(item).Error; err != nil {
		return nil, fmt.Errorf("failed to create {{.RepositoryName}}: %w", err)
	}
	return item, nil
}

// CreateMany adds multiple {{.RepositoryName}} entities to the database
func (r *{{.RepositoryName}}Repo) CreateMany(items []*models.{{.RepositoryName}}) ([]*models.{{.RepositoryName}}, error) {
	if err := r.DB.Create(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to create multiple {{.RepositoryName}}: %w", err)
	}
	return items, nil
}

// FindByID retrieves a {{.RepositoryName}} by its ID
func (r *{{.RepositoryName}}Repo) FindByID(id uint) (*models.{{.RepositoryName}}, error) {
	var item models.{{.RepositoryName}}
	if err := r.DB.First(&item, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find {{.RepositoryName}} by ID: %w", err)
	}
	return &item, nil
}

// GetAll retrieves all {{.RepositoryName}} entities from the database
func (r *{{.RepositoryName}}Repo) GetAll() ([]*models.{{.RepositoryName}}, error) {
	var items []*models.{{.RepositoryName}}
	if err := r.DB.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get all {{.RepositoryName}}: %w", err)
	}
	return items, nil
}

// GetPaged retrieves a paginated list of {{.RepositoryName}} entities
func (r *{{.RepositoryName}}Repo) GetPaged(page, pageSize int) ([]*models.{{.RepositoryName}}, error) {
	var items []*models.{{.RepositoryName}}
	offset := (page - 1) * pageSize
	if err := r.DB.Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get paginated {{.RepositoryName}}: %w", err)
	}
	return items, nil
}

// GetAllSorted retrieves all {{.RepositoryName}} entities sorted by a field
func (r *{{.RepositoryName}}Repo) GetAllSorted(orderBy string, ascending bool) ([]*models.{{.RepositoryName}}, error) {
	var items []*models.{{.RepositoryName}}
	if ascending {
		if err := r.DB.Order(orderBy).Find(&items).Error; err != nil {
			return nil, fmt.Errorf("failed to get sorted {{.RepositoryName}}: %w", err)
		}
	} else {
		if err := r.DB.Order(orderBy + " desc").Find(&items).Error; err != nil {
			return nil, fmt.Errorf("failed to get sorted {{.RepositoryName}}: %w", err)
		}
	}
	return items, nil
}

// UpdateOne updates a single {{.RepositoryName}} entity
func (r *{{.RepositoryName}}Repo) UpdateOne(item *models.{{.RepositoryName}}) (*models.{{.RepositoryName}}, error) {
	if err := r.DB.Save(item).Error; err != nil {
		return nil, fmt.Errorf("failed to update {{.RepositoryName}}: %w", err)
	}
	return item, nil
}

// UpdateByID updates a {{.RepositoryName}} entity by its ID
func (r *{{.RepositoryName}}Repo) UpdateByID(id uint, item *models.{{.RepositoryName}}) (*models.{{.RepositoryName}}, error) {
	if err := r.DB.Model(&models.{{.RepositoryName}}{}).Where("id = ?", id).Updates(item).Error; err != nil {
		return nil, fmt.Errorf("failed to update {{.RepositoryName}} by ID: %w", err)
	}
	return item, nil
}

// UpdateByEntity updates a {{.RepositoryName}} entity based on the entity itself
func (r *{{.RepositoryName}}Repo) UpdateByEntity(item *models.{{.RepositoryName}}) (*models.{{.RepositoryName}}, error) {
	if err := r.DB.Save(item).Error; err != nil {
		return nil, fmt.Errorf("failed to update {{.RepositoryName}} by entity: %w", err)
	}
	return item, nil
}

// UpdateMany updates multiple {{.RepositoryName}} entities
func (r *{{.RepositoryName}}Repo) UpdateMany(items []*models.{{.RepositoryName}}) ([]*models.{{.RepositoryName}}, error) {
	if err := r.DB.Save(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to update multiple {{.RepositoryName}}: %w", err)
	}
	return items, nil
}

// SoftDelete marks a {{.RepositoryName}} entity as deleted (without actually deleting it)
func (r *{{.RepositoryName}}Repo) SoftDelete(id uint) error {
	if err := r.DB.Model(&models.{{.RepositoryName}}{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
		return fmt.Errorf("failed to soft delete {{.RepositoryName}}: %w", err)
	}
	return nil
}

// DeleteOne deletes a single {{.RepositoryName}} by its ID
func (r *{{.RepositoryName}}Repo) DeleteOne(id uint) error {
	if err := r.DB.Delete(&models.{{.RepositoryName}}{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete {{.RepositoryName}} by ID: %w", err)
	}
	return nil
}

// DeleteMany deletes multiple {{.RepositoryName}} entities by their IDs
func (r *{{.RepositoryName}}Repo) DeleteMany(ids []uint) error {
	if err := r.DB.Where("id IN ?", ids).Delete(&models.{{.RepositoryName}}{}).Error; err != nil {
		return fmt.Errorf("failed to delete multiple {{.RepositoryName}} by IDs: %w", err)
	}
	return nil
}

// Count returns the total number of {{.RepositoryName}} records in the database
func (r *{{.RepositoryName}}Repo) Count() (int64, error) {
	var count int64
	if err := r.DB.Model(&models.{{.RepositoryName}}{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count {{.RepositoryName}}: %w", err)
	}
	return count, nil
}
`
