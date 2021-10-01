package output

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"edholm.dev/envconfig-generate/internal/tagparser"
)

const markdownTemplate = `# Available configuration
{{- range . }} 
## {{ .Package }}/{{ .Name }} 
{{- range .Options }}
` + "`" + `{{- .Name }}
	{{- if .Default }}={{ .Default }}
	{{- end }}
	{{- if .Required }}=<required-to-be-set>
	{{- end }}` + "`" + `
{{- end }}
{{ end }}
`

func ToMarkdown(ctx context.Context, options []tagparser.AvailableConfig) ([]byte, error) {
	tmpl := template.New("markdown")
	parse, err := tmpl.Parse(markdownTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	bb := new(bytes.Buffer)
	if err = parse.Execute(bb, options); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	return bb.Bytes(), nil
}
