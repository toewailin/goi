package templates

// DTOTemplate - Template for generating a DTO struct used for incoming requests.
const DTOTemplate = `package dto

// {{.DtoName}} represents the data structure for an incoming client request.
type {{.DtoName}} struct {
    Name        string    ` + "`" + `json:"name" binding:"required,min=3,max=100"` + "`" + ` // Example: Required field with validation
    Description string    ` + "`" + `json:"description,omitempty"` + "`" + ` // Example: Optional field
    IsActive    bool      ` + "`" + `json:"is_active"` + "`" + ` // Example: Boolean field
    Count       int       ` + "`" + `json:"count" binding:"gte=0"` + "`" + ` // Example: Integer with min value validation
{{- range .Fields}}
    {{.FieldName}} {{.FieldType}} ` + "`" + `json:"{{.JsonTag}}"` + "`" + `{{if .BindingTag}} ` + "`" + `binding:"{{.BindingTag}}"` + "`" + `{{end}}
{{- end}}
}
`
