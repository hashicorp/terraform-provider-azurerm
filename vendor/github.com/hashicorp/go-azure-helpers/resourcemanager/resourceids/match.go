// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourceids

import (
	"reflect"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/features"
)

// Match compares two instances of the same ResourceId and determines whether they are a match
//
// Whilst it might seem fine to compare the result of the `.ID()` function, that doesn't account
// for Resource ID Segments which need to be compared as case-insensitive.
//
// As such whilst this function is NOT exposing that functionality right now, it will when the
// centralised feature-flag for this is rolled out.
func Match(first, second ResourceId) bool {
	// since we're comparing interface types, ensure the two underlying types are the same
	if reflect.TypeOf(first) != reflect.TypeOf(second) {
		return false
	}

	parser := NewParserFromResourceIdType(first)
	firstParsed, err := parser.Parse(first.ID(), true)
	if err != nil {
		return false
	}
	secondParsed, err := parser.Parse(second.ID(), true)
	if err != nil {
		return false
	}
	firstVal := firstParsed.Parsed
	secondVal := secondParsed.Parsed
	if len(firstVal) != len(secondVal) {
		return false
	}
	for key, val := range firstVal {
		otherVal, ok := secondVal[key]
		if !ok {
			return false
		}

		segment := parser.namedSegment(key)
		if segment == nil {
			return false
		}

		if features.TreatUserSpecifiedSegmentsAsCaseInsensitive && segment.Type == UserSpecifiedSegmentType {
			if !strings.EqualFold(val, otherVal) {
				return false
			}
		} else if val != otherVal {
			return false
		}
	}

	return true
}
