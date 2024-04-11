// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourceids

import (
	"fmt"
	"reflect"
	"strings"
)

var _ error = SegmentNotSpecifiedError{}

type SegmentNotSpecifiedError struct {
	parseResult    ParseResult
	resourceId     ResourceId
	resourceIdName string
	segmentName    string
}

// NewSegmentNotSpecifiedError returns a SegmentNotSpecifiedError for the provided Resource ID, segment and parseResult combination
func NewSegmentNotSpecifiedError(id ResourceId, segmentName string, parseResult ParseResult) SegmentNotSpecifiedError {
	// Resource ID types must be in the format {Name}Id
	resourceIdTypeName := reflect.ValueOf(id).Type().Name()
	if resourceIdTypeName == "" {
		resourceIdTypeName = reflect.ValueOf(id).Elem().Type().Name()
	}
	resourceIdName := strings.TrimSuffix(resourceIdTypeName, "Id")
	return SegmentNotSpecifiedError{
		resourceIdName: resourceIdName,
		resourceId:     id,
		segmentName:    segmentName,
		parseResult:    parseResult,
	}
}

// Error returns a detailed error message highlighting the issues found when parsing this Resource ID Segment.
func (e SegmentNotSpecifiedError) Error() string {
	expectedId := buildExpectedResourceId(e.resourceId.Segments())
	position := findPositionOfSegment(e.segmentName, e.resourceId.Segments())
	if position == nil {
		return fmt.Sprintf("internal-error: couldn't determine the position for segment %q", e.segmentName)
	}
	description, err := descriptionForSegment(e.segmentName, e.resourceId.Segments())
	if err != nil {
		return fmt.Sprintf("internal-error: building description for segment: %+v", err)
	}

	return fmt.Sprintf(`parsing the %[1]s ID: the segment at position %[2]d didn't match

Expected a %[1]s ID that matched:

> %[3]s

However this value was provided:

> %[4]s

The parsed Resource ID was missing a value for the segment at position %[2]d
(which %[5]s).

`, e.resourceIdName, *position, expectedId, e.parseResult.RawInput, *description)
}
