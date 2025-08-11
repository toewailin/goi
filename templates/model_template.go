package templates

// model_template.go - Template for generating models
const ModelTemplate = `package models

import "gorm.io/gorm"

// {{.ModelName}} represents the {{.ModelName}} model in the database
type {{.ModelName}} struct {
	gorm.Model
	// Define fields here
}
`
