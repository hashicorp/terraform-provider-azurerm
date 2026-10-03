// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/markdown"
)

// S013 validates that the `id` attribute is at the top of Attributes Reference section
type S013 struct{}

var _ Rule = new(S013)

func (s S013) ID() string   { return "S013" }
func (s S013) Name() string { return "Attributes ID at Top" }
func (s S013) Description() string {
	return "Validates that the `id` attribute is at the top of Attributes Reference section"
}

// Run validates that id attribute is at top. Fix is not supported because moving lines
// would invalidate line numbers used by other rules.
func (s S013) Run(d *data.TerraformNodeData, _ bool) []error {
	if SkipRule(d.Type, d.Name, s.ID()) || !d.Document.Exists {
		return nil
	}

	var section *markdown.AttributesSection
	for _, sec := range d.Document.Sections {
		if attrSec, ok := sec.(*markdown.AttributesSection); ok {
			section = attrSec
			break
		}
	}
	if section == nil {
		return nil
	}

	content := section.GetContent()
	firstAttrIdx, idAttrIdx := -1, -1

	for idx, line := range content {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)

		// Stop at block definition markers:
		// - "---"
		// - "A/An `xxx` block ..."
		// - "The `xxx` block ..."
		if trimmed == "---" ||
			strings.HasPrefix(lower, "a `") ||
			strings.HasPrefix(lower, "an `") ||
			strings.HasPrefix(lower, "the `") {
			break
		}

		// Only check top-level attributes (lines starting with "* `")
		if strings.HasPrefix(trimmed, "* `") {
			if firstAttrIdx == -1 {
				firstAttrIdx = idx
			}
			if strings.HasPrefix(trimmed, "* `id`") {
				idAttrIdx = idx
				break
			}
		}
	}

	if idAttrIdx != -1 && idAttrIdx != firstAttrIdx {
		firstLine := content[firstAttrIdx]
		idLine := content[idAttrIdx]

		return []error{NewValidationIssue(s.ID(), s.Name(), "id",
			"`id` should be the first attribute",
			d.Document.Path, firstLine, idLine)}
	}

	return nil
}
