// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package commands

import (
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/helpers"
	"github.com/iancoleman/strcase"
)

var config = &helpers.Configuration{}

var TplFuncMap = template.FuncMap{
	"ToLower":                       strings.ToLower,
	"ToTitle":                       strings.Title,
	"ToCamel":                       strcase.ToCamel,
	"ToSnake":                       strcase.ToSnake,
	"TfName":                        helpers.TerraformResourceName,
	"ToString":                      helpers.ToString,
	"ToDelim":                       strcase.ToDelimited,
	"ToDelimTitle":                  helpers.ToDelimTitle,
	"PrefixedDescriptionString":     helpers.PrefixedDescriptionString,
	"PrefixedLabelString":           helpers.PrefixedLabelString,
	"SchemaItemFormatter":           helpers.SchemaItemFormatter,
	"SchemaItemFormatterAttributes": helpers.SchemaItemFormatterAttributes,
	"SchemaItemFormatterSpecial":    helpers.SchemaItemFormatterSpecial,
}

func init() {
	config = helpers.LoadConfig()
}
