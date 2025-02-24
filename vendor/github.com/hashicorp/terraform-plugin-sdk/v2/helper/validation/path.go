// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"github.com/hashicorp/go-cty/cty"
)

// PathMatches compares two Paths for equality. For cty.IndexStep,
// unknown key values are treated as an Any qualifier and will
// match any index step of the same type.
func PathMatches(p cty.Path, other cty.Path) bool {
	if len(p) != len(other) {
		return false
	}

	for i := range p {
		pv := p[i]
		switch pv := pv.(type) {
		case cty.GetAttrStep:
			ov, ok := other[i].(cty.GetAttrStep)
			if !ok || pv != ov {
				return false
			}
		case cty.IndexStep:
			ov, ok := other[i].(cty.IndexStep)
			if !ok {
				return false
			}

			// Sets need special handling since their Type is the entire object
			// with attributes.
			if pv.Key.Type().IsObjectType() && ov.Key.Type().IsObjectType() {
				if !pv.Key.IsKnown() || !ov.Key.IsKnown() {
					break
				}
			}
			if !pv.Key.Type().Equals(ov.Key.Type()) {
				return false
			}

			if pv.Key.IsKnown() && ov.Key.IsKnown() {
				if !pv.Key.RawEquals(ov.Key) {
					return false
				}
			}
		default:
			// Any invalid steps default to evaluating false.
			return false
		}
	}

	return true
}
