package rule

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/util"
)

// S007 validates and fixes Required/Optional markers in documentation
type S007 struct{}

var _ Rule = new(S007)

func (s S007) ID() string {
	return "S007"
}

func (s S007) Name() string {
	return "Properties Requiredness in Document"
}

func (s S007) Description() string {
	return "Determines whether all properties are correctly marked as Required or Optional in documentation"
}

func (s S007) Run(d *data.TerraformNodeData, fix bool) []error {
	if !d.Document.Exists {
		return nil
	}

	if d.Type == data.ResourceTypeData {
		return nil
	}

	if d.SchemaProperties == nil || d.DocumentArguments == nil {
		return nil
	}

	return forEachSchemaProperty(s.ID(), d, "", d.SchemaProperties, d.DocumentArguments, d.DocumentArguments.BlockDefinitions,
		func(fullPath string, schemaProperty *models.SchemaProperty, docProperty *models.DocumentProperty) error {
			return s.checkPropertyRequiredness(d, fullPath, schemaProperty, docProperty, fix)
		},
	)
}

func (s S007) checkPropertyRequiredness(
	d *data.TerraformNodeData,
	fullPath string,
	schemaProperty *models.SchemaProperty,
	docProperty *models.DocumentProperty,
	fix bool,
) error {
	if schemaProperty.Required && !docProperty.Required {
		expected := s.replaceRequiredness(docProperty.Content, "(Optional)", "(Required)")
		issue := NewValidationIssue(
			s.ID(),
			s.Name(),
			fullPath,
			fmt.Sprintf("`%s` should be marked as Required", util.Bold(fullPath)),
			d.Document.Path,
			docProperty.Content,
			expected,
		)
		if fix {
			s.applyFix(d, docProperty, true)
		}
		return issue
	}

	if schemaProperty.Optional && !docProperty.Optional {
		expected := s.replaceRequiredness(docProperty.Content, "(Required)", "(Optional)")
		issue := NewValidationIssue(
			s.ID(),
			s.Name(),
			fullPath,
			fmt.Sprintf("`%s` should be marked as Optional", util.Bold(fullPath)),
			d.Document.Path,
			docProperty.Content,
			expected,
		)
		if fix {
			s.applyFix(d, docProperty, false)
		}
		return issue
	}

	return nil
}

// replaceRequiredness replaces one requiredness marker with another
func (s S007) replaceRequiredness(line, from, to string) string {
	if strings.Contains(line, from) {
		return strings.Replace(line, from, to, 1)
	} else {
		// add after the first -
		if idx := strings.Index(line, " - "); idx > 0 {
			line = line[:idx+3] + to + " " + line[idx+3:]
		} else {
			// no dash add after second `
			idx = strings.Index(line, "`")
			idx += strings.Index(line[idx+1:], "`") + 1
			line = line[:idx+1] + " " + to + line[idx+1:]
		}
	}
	return line
}

func (s S007) applyFix(d *data.TerraformNodeData, docProperty *models.DocumentProperty, shouldBeRequired bool) {
	if d.Document == nil {
		return
	}

	argsSection := d.Document.GetArgumentsSection()
	if argsSection == nil {
		return
	}

	content := argsSection.GetContent()
	lineIdx := docProperty.Line
	from := "(Required)"
	to := "(Optional)"
	if shouldBeRequired {
		from = "(Optional)"
		to = "(Required)"
	}

	if lineIdx >= 0 && lineIdx < len(content) {
		line := content[lineIdx]
		fixedLine := s.replaceRequiredness(line, from, to)
		content[lineIdx] = fixedLine
		argsSection.SetContent(content)
		d.Document.HasChange = true
	}
}
