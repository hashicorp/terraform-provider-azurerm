// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import (
	"fmt"
	"strings"
)

var (
	_ PropertyRule = redundantIsPrefix{}
	_ Fixer        = redundantIsPrefix{}
)

// redundantIsPrefix flags boolean properties that begin with a redundant "is_"
// verb (e.g. is_storage_enabled should be storage_enabled).
type redundantIsPrefix struct{}

func (redundantIsPrefix) ID() string   { return "SL015" }
func (redundantIsPrefix) Name() string { return "redundant-is-prefix" }

func (redundantIsPrefix) Description() string {
	return `Boolean properties should not begin with a redundant "is_" verb.`
}

func (redundantIsPrefix) DefaultSeverity() Severity { return SeverityWarning }

func (redundantIsPrefix) FixHint() string { return `drop the "is_" prefix` }

func (r redundantIsPrefix) CheckProperty(ctx PropertyContext) []Finding {
	if ctx.Schema.Type != TypeBool {
		return nil
	}
	if !strings.HasPrefix(ctx.Name, "is_") {
		return nil
	}

	preferred := strings.TrimPrefix(ctx.Name, "is_")
	if preferred == "" {
		return nil
	}

	return []Finding{propertyFindingFix(r, ctx,
		fmt.Sprintf("boolean property %q has a redundant \"is_\" prefix (%q)", ctx.Path, preferred),
		fmt.Sprintf("rename %q to %q", ctx.Name, preferred),
	)}
}
