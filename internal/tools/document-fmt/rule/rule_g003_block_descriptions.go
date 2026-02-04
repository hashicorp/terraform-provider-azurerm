// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/differror"
)

type G003 struct{}

var (
	_ Rule = G003{}

	// Matches block property descriptions like:
	// * `name` - A/An/One or more/A list of `name` block(s) as defined/documented below/above.
	// With optional trailing description after the period
	blockDescriptionRegex = regexp.MustCompile(
		`^\* ` + "`" + `([^` + "`" + `]+)` + "`" + ` - ` + // property name
			`(A |An |One or more |A list of )?` + // optional quantifier
			"`" + `([^` + "`" + `]+)` + "`" + ` ` + // block name
			`blocks? as (defined|documented) (below|above)` + // defined/documented below/above
			`(.*)$`) // optional trailing text
)

// schemaProperty contains information about a Terraform schema property
type schemaProperty struct {
	Name     string
	IsBlock  bool
	MinItems int
	MaxItems int
	IsSet    bool // TypeSet vs TypeList
	Required bool
	Optional bool
	Computed bool
}

func (r G003) ID() string {
	return "G003"
}

func (r G003) Name() string {
	return "Block Descriptions"
}

func (r G003) Description() string {
	return "validates and fixes block property descriptions using schema data for accuracy"
}

// getSchemaProperties extracts block properties from the schema recursively
func getSchemaProperties(s map[string]*schema.Schema, prefix string) map[string]schemaProperty {
	props := make(map[string]schemaProperty)

	for name, prop := range s {
		fullName := name
		if prefix != "" {
			fullName = prefix + "." + name
		}

		sp := schemaProperty{
			Name:     name,
			MinItems: prop.MinItems,
			MaxItems: prop.MaxItems,
			IsSet:    prop.Type == schema.TypeSet,
			Required: prop.Required,
			Optional: prop.Optional,
			Computed: prop.Computed,
		}

		// Check if this is a block (has nested schema)
		if prop.Elem != nil {
			if nestedResource, ok := prop.Elem.(*schema.Resource); ok {
				sp.IsBlock = true
				props[name] = sp

				// Recursively get nested properties
				for k, v := range getSchemaProperties(nestedResource.Schema, fullName) {
					props[fullName+"."+k] = v
				}
				continue
			}
		}

		props[name] = sp
	}

	return props
}

func (r G003) Run(d *data.TerraformNodeData, fix bool) []error {
	errs := make([]error, 0)

	// Build schema property map if available
	schemaProps := make(map[string]schemaProperty)
	if d.Resource != nil && d.Resource.Schema != nil {
		schemaProps = getSchemaProperties(d.Resource.Schema, "")
	}

	for _, section := range d.Document.Sections {
		content := section.GetContent()
		for idx, line := range content {
			matches := blockDescriptionRegex.FindStringSubmatch(line)
			if matches == nil {
				continue
			}

			propName := matches[1]
			quantifier := matches[2]
			blockName := matches[3]
			verb := matches[4]      // "defined" or "documented"
			direction := matches[5] // "below" or "above"
			trailing := matches[6]  // optional trailing text

			hasErrors := false
			var errorMessages []string

			// Look up schema info for this property
			schemaProp, hasSchemaInfo := schemaProps[propName]

			// Check 1: "documented" should be "defined"
			if verb == "documented" {
				hasErrors = true
				verb = "defined"
				errorMessages = append(errorMessages, "'documented' should be 'defined'")
			}

			// Determine expected quantifier based on schema
			expectedQuantifier := quantifier
			expectedBlockWord := "block"

			if hasSchemaInfo && schemaProp.IsBlock {
				// Use schema to determine singular vs plural
				if schemaProp.MaxItems == 1 {
					// Single block: "A" or "An"
					firstChar := strings.ToLower(string(blockName[0]))
					if strings.ContainsAny(firstChar, "aeiou") {
						expectedQuantifier = "An "
					} else {
						expectedQuantifier = "A "
					}
					expectedBlockWord = "block"
				} else {
					// Multiple blocks allowed
					if schemaProp.MinItems >= 1 {
						expectedQuantifier = "One or more "
					} else {
						expectedQuantifier = "A list of "
					}
					expectedBlockWord = "blocks"
				}
			} else if quantifier != "" {
				// No schema info, fall back to regex-based validation
				if quantifier == "A " || quantifier == "An " {
					firstChar := strings.ToLower(string(blockName[0]))
					if strings.ContainsAny(firstChar, "aeiou") {
						expectedQuantifier = "An "
					} else {
						expectedQuantifier = "A "
					}
					expectedBlockWord = "block"
				} else if quantifier == "One or more " || quantifier == "A list of " {
					expectedBlockWord = "blocks"
				}
			}

			// Check if quantifier needs fixing
			if quantifier != expectedQuantifier && expectedQuantifier != "" {
				hasErrors = true
				if hasSchemaInfo {
					if schemaProp.MaxItems == 1 {
						errorMessages = append(errorMessages, fmt.Sprintf("schema indicates MaxItems=1, should use '%s'", strings.TrimSpace(expectedQuantifier)))
					} else {
						errorMessages = append(errorMessages, fmt.Sprintf("schema indicates multiple items allowed, should use '%s'", strings.TrimSpace(expectedQuantifier)))
					}
				} else {
					errorMessages = append(errorMessages, fmt.Sprintf("'%s' should be '%s'", strings.TrimSpace(quantifier), strings.TrimSpace(expectedQuantifier)))
				}
			}

			if hasErrors {
				current := line
				expected := fmt.Sprintf("* `%s` - %s`%s` %s as %s %s%s",
					propName,
					expectedQuantifier,
					blockName,
					expectedBlockWord,
					verb,
					direction,
					trailing,
				)

				errs = append(errs, differror.New(
					fmt.Sprintf("%s: Block description issues: %s", IdAndName(r), strings.Join(errorMessages, ", ")),
					current,
					expected,
				))

				if fix {
					d.Document.HasChange = true
					content[idx] = expected
				}
			}
		}

		if fix && d.Document.HasChange {
			section.SetContent(content)
		}
	}

	return errs
}
