package rule

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
)

type S006 struct{}

var _ Rule = new(S006)

func (s S006) ID() string {
	return "S006"
}

func (s S006) Name() string {
	return "Arguments Exist in Document"
}

func (s S006) Description() string {
	return "Determines whether all arguments defined in schema are documented"
}

func (s S006) Run(d *data.TerraformNodeData, fix bool) []error {
	var errs []error

	if !d.Document.Exists {
		return errs
	}

	// If either is nil, we can't properly check, so short-circuit
	// TODO: rework so this is not needed and we can just rely on the checks and loops below
	if d.SchemaProperties == nil || d.DocumentArguments == nil {
		return errs
	}

	errs = append(errs, s.check(d.SchemaProperties, d.DocumentArguments, "")...)

	if errs != nil {
	}

	return errs
}

func (s S006) check(schema *data.Properties, documentation *data.Properties, parentArgument string) []error {
	errs := make([]error, 0)
	for name, property := range schema.Objects {
		if !property.Optional && !property.Required {
			continue
		}

		fullName := name
		if parentArgument != "" {
			fullName = parentArgument + "." + name
		}

		// rename to something else
		split := strings.Split(name, ".")
		lastChild := split[len(split)-1]

		documentedProperty, ok := documentation.Objects[lastChild]
		if !ok {
			errs = append(errs, fmt.Errorf("%s: expected argument `%s` to be present in the documentation", IdAndName(s), fullName))
			continue
		}

		if property.Nested != nil {
			if documentedProperty.Nested == nil {
				// unlikely
				errs = append(errs, fmt.Errorf("%s: expected `%s` to have documented nested arguments", IdAndName(s), fullName))
				continue
			}
			errs = append(errs, s.check(property.Nested, documentedProperty.Nested, fullName)...)
		}
	}

	return errs
}
