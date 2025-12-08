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
// S<number> -- section specific rules, e.g. API section validatio
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
	S006{}.ID(): S006{}, // Arguments Format
	S007{}.ID(): S007{}, // Arguments Requiredness
	S008{}.ID(): S008{}, // Arguments ForceNew
	S009{}.ID(): S009{}, // Arguments Existence
	S010{}.ID(): S010{}, // Arguments Default Value
	S011{}.ID(): S011{}, // Arguments Possible Values
}

// RulesRequiringPackageLoad contains rule IDs that need package loading and API analysis
// These rules require expensive package loading and SSA analysis
var RulesRequiringPackageLoad = map[string]bool{
	S001{}.ID(): true, // API Section validation needs API data from loaded packages
}

// RulesRequiringPackageLoad contains rule IDs that need normalize markdowns
var RulesRequiringNormalizingMd = map[string]bool{
	S006{}.ID(): true,
	S007{}.ID(): true,
	S008{}.ID(): true,
	S009{}.ID(): true,
	S010{}.ID(): true,
	S011{}.ID(): true,
}

func ShouldLoadPackages(ruleIDs []string) bool {
	for _, ruleID := range ruleIDs {
		if RulesRequiringPackageLoad[ruleID] {
			return true
		}
	}
	return false
}

func ShouldNormalizeMd(ruleIDs []string) bool {
	for _, ruleID := range ruleIDs {
		if RulesRequiringNormalizingMd[ruleID] {
			return true
		}
	}
	return false
}
