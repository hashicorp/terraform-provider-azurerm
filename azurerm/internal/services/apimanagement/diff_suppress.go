package apimanagement

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
)

// XmlWithDotNetInterpolationsDiffSuppress is a Diff Suppress Func for when the XML contains
// .net interpolations, and thus isn't valid XML to parse
// whilst really we should be parsing the XML Tokens and skipping over the error - in practice
func XmlWithDotNetInterpolationsDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// try parsing this as valid xml if we can, to handle ordering differences
	same := suppress.XmlDiff(k, old, new, d)
	if same {
		return same
	}

	// otherwise best-effort this via string comparison
	oldVal := normalizeXmlWithDotNetInterpolationsString(old)
	newVal := normalizeXmlWithDotNetInterpolationsString(new)
	return oldVal == newVal
}

// normalizeXmlWithDotNetInterpolationsString is intended as a fallback to diff two xml strings
// containing .net interpolations, which means that they aren't directly valid xml
// whilst we /could/ xml.EscapeString these that encodes the entire string, rather than the expression
// we could do that as a potential extension, but this seems sufficient in testing :shrug:
func normalizeXmlWithDotNetInterpolationsString(input string) string {
	value := input

	value = strings.ReplaceAll(value, "\n", "")
	value = strings.ReplaceAll(value, "\r", "")
	value = strings.ReplaceAll(value, "\t", "")
	value = strings.ReplaceAll(value, "    ", "")
	value = strings.ReplaceAll(value, "  ", "")
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "&quot;", "\"")
	value = strings.ReplaceAll(value, "&gt;", ">")
	value = strings.ReplaceAll(value, "&lt;", "<")
	value = strings.ReplaceAll(value, "&amp;", "&")
	value = strings.ReplaceAll(value, "&apos;", "'")

	return value
}
