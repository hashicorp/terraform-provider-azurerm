package generators

import (
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/generator-tests/helpers"
	"github.com/iancoleman/strcase"
)

var TplFuncMap = template.FuncMap{
	"ToLower":                   strings.ToLower,
	"ToTitle":                   strings.Title,
	"ToCamel":                   strcase.ToCamel,
	"ToSnake":                   strcase.ToSnake,
	"TfName":                    helpers.TerraformResourceName,
	"ToString":                  helpers.ToString,
	"ToDelim":                   strcase.ToDelimited,
	"ToDelimTitle":              helpers.ToDelimTitle,
	"PrefixedDescriptionString": helpers.PrefixedDescriptionString,
	"PrefixedLabelString":       helpers.PrefixedLabelString,
}
