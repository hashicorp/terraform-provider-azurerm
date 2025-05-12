package template

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/data"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Render(data *data.ResourceData, text string) ([]string, error) {
	var err error
	var b bytes.Buffer

	tmpl := template.New("template")
	tmpl.Funcs(map[string]interface{}{
		"lower": strings.ToLower,
		"title": toTitle,
	})

	tmpl, err = tmpl.Parse(text)
	if err != nil {
		return nil, err
	}

	err = tmpl.Execute(&b, data)

	content := strings.ReplaceAll(b.String(), "[bt]", "`")
	return strings.Split(content, "\n"), nil
}

func toTitle(s string) string {
	caser := cases.Title(language.English)
	return caser.String(s)
}
