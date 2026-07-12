// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package rules

import "fmt"

var _ PropertyRule = collectionElemRequired{}

// collectionElemRequired flags a list or set property that does not define an
// Elem, which the terraform-plugin-sdk requires.
type collectionElemRequired struct{}

func (collectionElemRequired) ID() string   { return "SL005" }
func (collectionElemRequired) Name() string { return "collection-elem-required" }

func (collectionElemRequired) Description() string {
	return "A TypeList or TypeSet property must define an Elem."
}

func (collectionElemRequired) DefaultSeverity() Severity { return SeverityError }

func (r collectionElemRequired) CheckProperty(ctx PropertyContext) []Finding {
	if IsCollection(ctx.Schema) && ctx.Schema.Elem == nil {
		return []Finding{propertyFinding(r, ctx, fmt.Sprintf("collection property %q (%s) does not define an Elem", ctx.Path, ctx.Schema.Type))}
	}

	return nil
}
