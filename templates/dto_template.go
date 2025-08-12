package templates

// DTOTemplate - Template for generating a DTO struct
const DTOtemplate = `package {{.PackageName}}

type {{.StructName}} struct {
{{- range .Fields}}
    {{.FieldName}} {{.FieldType}} ` + "`" + `json:"{{.JsonTag}}"` + "`" + `{{if .BindingTag}} ` + "`" + `binding:"{{.BindingTag}}"` + "`" + `{{end}}
{{- end}}
}
`