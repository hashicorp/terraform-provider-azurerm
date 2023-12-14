// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourceids

import (
	"fmt"
	"strings"
)

// descriptionForSegments returns a summary/description for the provided Resource ID Segments
func descriptionForSegments(input []Segment) (*string, error) {
	lines := make([]string, 0)

	for i, segment := range input {
		description, err := descriptionForSpecifiedSegment(segment)
		if err != nil {
			return nil, fmt.Errorf("building description for segment %q: %+v", segment.Name, err)
		}

		lines = append(lines, fmt.Sprintf("* Segment %d - this %s", i, *description))
	}

	out := strings.Join(lines, "\n")
	return &out, nil
}

// descriptionForSegment returns a friendly description for the Segment specified in segmentName
func descriptionForSegment(segmentName string, segments []Segment) (*string, error) {
	for _, segment := range segments {
		if segment.Name != segmentName {
			continue
		}

		return descriptionForSpecifiedSegment(segment)
	}

	return nil, fmt.Errorf("the segment %q was not defined for this Resource ID", segmentName)
}

// descriptionForSpecifiedSegment returns a friendly description for the Segment
func descriptionForSpecifiedSegment(segment Segment) (*string, error) {
	// NOTE: do not use round brackets within these error messages, since this description can be contained within one
	// the description will also be prefixed with a `which `
	switch segment.Type {
	case ConstantSegmentType:
		{
			if segment.PossibleValues == nil {
				return nil, fmt.Errorf("the Segment %q defined a Constant with no PossibleValues", segment.Name)
			}

			// intentionally format these as quoted string values
			values := make([]string, 0)
			for _, v := range *segment.PossibleValues {
				values = append(values, fmt.Sprintf("%q", v))
			}

			msg := fmt.Sprintf("should be a Constant with one of the following values [%s]", strings.Join(values, ", "))
			return &msg, nil
		}

	case ResourceGroupSegmentType:
		{
			msg := "should be the name of the Resource Group"
			return &msg, nil
		}

	case ResourceProviderSegmentType:
		{
			msg := fmt.Sprintf("should be the name of the Resource Provider [for example '%s']", segment.ExampleValue)
			return &msg, nil
		}

	case ScopeSegmentType:
		{
			msg := fmt.Sprintf("specifies the Resource ID that should be used as a Scope [for example '%s']", segment.ExampleValue)
			return &msg, nil
		}

	case StaticSegmentType:
		{
			if segment.FixedValue == nil {
				return nil, fmt.Errorf("the Segment %q defined a Static Segment with no FixedValue", segment.Name)
			}
			msg := fmt.Sprintf("should be the literal value %q", *segment.FixedValue)
			return &msg, nil
		}

	case SubscriptionIdSegmentType:
		{
			msg := "should be the UUID of the Azure Subscription"
			return &msg, nil
		}

	case UserSpecifiedSegmentType:
		{
			name := strings.TrimSuffix(segment.Name, "Name")
			msg := fmt.Sprintf("should be the user specified value for this %s [for example %q]", name, segment.ExampleValue)
			return &msg, nil
		}
	}

	return nil, fmt.Errorf("internal-error: the Segment Type %q was not implemented for Segment %q", string(segment.Type), segment.Name)
}

// buildExpectedResourceId iterates over the Resource ID to build up the "expected" value for the Resource ID
// this is done using the example segment values for each segment type.
func buildExpectedResourceId(segments []Segment) string {
	components := make([]string, 0)
	for _, v := range segments {
		components = append(components, strings.TrimPrefix(v.ExampleValue, "/"))
	}

	out := strings.Join(components, "/")
	return fmt.Sprintf("/%s", out)
}

// findPositionOfSegment returns the position of the segment specified by segmentName within segments
func findPositionOfSegment(segmentName string, segments []Segment) *int {
	for i, segment := range segments {
		if segment.Name == segmentName {
			return &i
		}
	}

	return nil
}

// summaryOfParsedSegments returns a summary of the parsed Resource ID Segments vs what we're expecting
func summaryOfParsedSegments(parsed ParseResult, segments []Segment) string {
	out := make([]string, 0)
	for i, v := range segments {
		val, ok := parsed.Parsed[v.Name]
		if !ok {
			out = append(out, fmt.Sprintf("* Segment %d - not found", i))
			continue
		}

		out = append(out, fmt.Sprintf("* Segment %d - parsed as %q", i, val))
	}

	return strings.Join(out, "\n")
}
