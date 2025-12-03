// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validator

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/rule"
)

// Validator.Run provides a separation between `data` and `rule` packages to avoid cyclic imports...
// TODO: review tool structure?

type Validator struct{}

// Rule execution order is important to avoid conflicts when fixing:
// Other rules can execute in any order
var fixRuleOrder = []string{
	rule.S006{}.ID(),
	rule.S007{}.ID(),
	rule.S008{}.ID(),
	rule.S009{}.ID(),
}

func (v Validator) Run(rules []string, data *data.TerraformNodeData, fix bool) {
	if fix {
		// In fix mode, execute rules in specific order to avoid conflicts
		v.runInOrder(rules, data, fix)
	} else {
		for _, id := range rules {
			if r := rule.GetRule(id); r != nil {
				errs := r.Run(data, fix)
				if len(errs) > 0 {
					data.Errors = append(data.Errors, errs...)
				}
			}
		}
	}
}

func (v Validator) runInOrder(rules []string, data *data.TerraformNodeData, fix bool) {
	// Create a set of requested rules for quick lookup
	requestedRules := make(map[string]bool)
	for _, id := range rules {
		requestedRules[id] = true
	}

	// Execute rules in priority order
	for _, id := range fixRuleOrder {
		if !requestedRules[id] {
			continue
		}
		if r := rule.GetRule(id); r != nil {
			errs := r.Run(data, fix)
			if len(errs) > 0 {
				data.Errors = append(data.Errors, errs...)
			}
		}
	}

	// Execute remaining rules that are not in priority list
	for _, id := range rules {
		if isInPriorityList(id) {
			continue
		}
		if r := rule.GetRule(id); r != nil {
			errs := r.Run(data, fix)
			if len(errs) > 0 {
				data.Errors = append(data.Errors, errs...)
			}
		}
	}
}

func isInPriorityList(id string) bool {
	for _, priorityID := range fixRuleOrder {
		if id == priorityID {
			return true
		}
	}
	return false
}
