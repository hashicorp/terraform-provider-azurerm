// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
)

type Rule interface {
	// ID returns the ID of the rule, e.g. `G001`.
	ID() string
	// Name returns a friendly name of the rule, e.g. `Validate API Section`.
	Name() string
	// Description returns an explanation of what this rule is validating.
	Description() string
	// Run implements the validation logic, the boolean parameter is true when cmd is "fix"
	Run(*data.TerraformNodeData, bool) []error
}

func GetRule(name string) Rule {
	if rule, ok := Registration[name]; ok {
		return rule
	}

	return nil
}

func IdAndName(r Rule) string {
	return fmt.Sprintf("%s - %s", r.ID(), r.Name())
}

// Registration for rules
// G<number> -- global/section agnostic rules, e.g. note formatting
// S<number> -- section specific rules, e.g. API section validation
var Registration = map[string]Rule{
	// global
	G001{}.ID(): G001{}, // Files Exist
	G002{}.ID(): G002{}, // Note Formatting

	// section
	S001{}.ID(): S001{}, // API Section
	S002{}.ID(): S002{}, // Timeouts Section
	S003{}.ID(): S003{}, // Title Section Heading
	S004{}.ID(): S004{}, // Arguments Section Heading
	S005{}.ID(): S005{}, // Attributes Section Heading
	S006{}.ID(): S006{}, // Arguments Exist in Document
	// S007{}.ID(): S007{}, // Attributes Exist in Document
}
