// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "fmt"

var _ PropertyRule = blockNeedsConstraint{}

// blockNeedsConstraint flags a user-settable block (a TypeList/TypeSet with a
// resource Elem) whose nested fields are all optional and none of which set
// AtLeastOneOf or ExactlyOneOf. Without a required field or conditional
// validation such a block can be configured empty.
type blockNeedsConstraint struct{}

func (blockNeedsConstraint) ID() string   { return "SL006" }
func (blockNeedsConstraint) Name() string { return "block-needs-constraint" }

func (blockNeedsConstraint) Description() string {
	return "A block with no required nested fields should set AtLeastOneOf/ExactlyOneOf on its optional fields so it cannot be configured empty."
}

func (blockNeedsConstraint) DefaultSeverity() Severity { return SeverityWarning }

func (r blockNeedsConstraint) CheckProperty(ctx PropertyContext) []Finding {
	if !ctx.Schema.Optional && !ctx.Schema.Required {
		return nil
	}

	blockSchema, ok := BlockElem(ctx.Schema)
	if !ok || len(blockSchema.Schema) == 0 {
		return nil
	}

	for _, child := range blockSchema.Schema {
		if child.Required {
			return nil
		}
		if len(child.AtLeastOneOf) > 0 || len(child.ExactlyOneOf) > 0 {
			return nil
		}
	}

	return []Finding{propertyFinding(r, ctx, fmt.Sprintf("block %q has no required fields and no AtLeastOneOf/ExactlyOneOf; add conditional validation so it cannot be configured empty", ctx.Path))}
}
