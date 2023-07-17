package logic

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed typed.gotpl
var tplString string

var (
	typedTpl = template.Must(template.New("typed").Parse(tplString))
)

func (m *MetaInfo) codeGen() ([]byte, error) {
	var buf bytes.Buffer
	if err := typedTpl.Execute(&buf, m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
