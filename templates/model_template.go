package templates

// model_template.go - Template for generating models
const ModelTemplate = `package models

import "gorm.io/gorm"

// {{.ModelName}} represents the {{.ModelName}} model in the database
type {{.ModelName}} struct {
    gorm.Model
    Name        string ` + "`" + `json:"name"` + "`" + `
    Description string ` + "`" + `json:"description,omitempty"` + "`" + `
    IsActive    bool   ` + "`" + `json:"is_active"` + "`" + `
    Count       int    ` + "`" + `json:"count"` + "`" + `
    // Define fields here
}
`
