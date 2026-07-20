// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/pr-guide-suggestions/schematree"
)

var sl013 = &Rule{
	ID:       "SL013",
	Name:     "id-reference-validation",
	Severity: Warning,
	Check:    checkSL013,
}

func checkSL013(res *schematree.Result, report ReportFunc) {
	for _, n := range res.All {
		if n.Schema == nil || !userSettable(n.Schema) {
			continue
		}
		if n.Schema.ValueType() != typeString || !strings.HasSuffix(n.Name, "_id") {
			continue
		}

		if !hasValidation(n.Schema) {
			report(n, fmt.Sprintf("ID reference %q has no validation; use a resource-specific ID validator (e.g. a commonids Validate...ID function) or validation.IsUUID", n.Path), nil)
			continue
		}
		if isSelectorFunc(n.Schema.FieldValue(fieldValidateFunc), "StringIsNotEmpty") {
			report(n, fmt.Sprintf("ID reference %q uses weak validation (accepts arbitrary strings); use a resource-specific ID validator or validation.IsUUID", n.Path), nil)
		}
	}
}
