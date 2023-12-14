// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourceids

import (
	"fmt"
	"reflect"
	"strings"
)

var _ error = NumberOfSegmentsDidntMatchError{}

type NumberOfSegmentsDidntMatchError struct {
	parseResult    ParseResult
	resourceId     ResourceId
	resourceIdName string
}

func NewNumberOfSegmentsDidntMatchError(id ResourceId, parseResult ParseResult) NumberOfSegmentsDidntMatchError {
	// Resource ID types must be in the format {Name}Id
	resourceIdName := strings.TrimSuffix(reflect.ValueOf(id).Type().Name(), "Id")
	return NumberOfSegmentsDidntMatchError{
		parseResult:    parseResult,
		resourceId:     id,
		resourceIdName: resourceIdName,
	}
}

// Error returns a detailed error message highlighting the issues found when parsing this Resource ID Segment.
func (e NumberOfSegmentsDidntMatchError) Error() string {
	expectedId := buildExpectedResourceId(e.resourceId.Segments())

	description, err := descriptionForSegments(e.resourceId.Segments())
	if err != nil {
		return fmt.Sprintf("internal-error: building description for segments: %+v", err)
	}

	parsedSegments := summaryOfParsedSegments(e.parseResult, e.resourceId.Segments())

	return fmt.Sprintf(`parsing the %[1]s ID: the number of segments didn't match

Expected a %[1]s ID that matched (containing %[2]d segments):

> %[3]s

However this value was provided (which was parsed into %[4]d segments):

> %[5]s

The following Segments are expected:

%[6]s

The following Segments were parsed:

%[7]s
`, e.resourceIdName, len(e.resourceId.Segments()), expectedId, len(e.parseResult.Parsed), e.parseResult.RawInput, *description, parsedSegments)
}
