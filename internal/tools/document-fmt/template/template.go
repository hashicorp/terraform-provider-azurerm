// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package template

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Render(data *data.TerraformNodeData, text string) ([]string, error) {
	var err error
	var b bytes.Buffer

	tmpl := template.New("template")
	tmpl.Funcs(map[string]interface{}{
		"lower":    strings.ToLower,
		"title":    toTitle,
		"id":       formattedID,
		"identity": formattedIdentity,
	})

	tmpl, err = tmpl.Parse(text)
	if err != nil {
		return nil, err
	}

	err = tmpl.Execute(&b, data)
	if err != nil {
		return nil, err
	}

	content := strings.ReplaceAll(b.String(), "[bt]", "`")
	return strings.Split(content, "\n"), nil
}

func toTitle(s string) string {
	caser := cases.Title(language.English)
	return caser.String(s)
}

func formattedID(id resourceids.ResourceId) string {
	result := "/"
	for _, segment := range id.Segments() {
		switch segment.Type {
		case resourceids.StaticSegmentType, resourceids.ResourceProviderSegmentType:
			// FixedValue _should_ be populated, but just in case
			if segment.FixedValue == nil {
				result += segment.ExampleValue
				break
			}

			result += *segment.FixedValue
		case resourceids.SubscriptionIdSegmentType:
			result += segment.ExampleValue
		case resourceids.ResourceGroupSegmentType, resourceids.UserSpecifiedSegmentType:
			result += "{" + segment.Name + "}"
		default:
			result += segment.Name
		}

		result += "/"
	}

	// could length check in loop but meh
	return strings.TrimSuffix(result, "/")
}

func formattedIdentity(id resourceids.ResourceId, idType pluginsdk.ResourceTypeForIdentity) string {
	result := "TODO Resource Identity Format"
	// - need to take into account identity "types" (e.g. virtual where we don't trim down to just `name` for the last field)
	// - need to ensure even spacing (ie `terraform fmt` output)

	return result
}
