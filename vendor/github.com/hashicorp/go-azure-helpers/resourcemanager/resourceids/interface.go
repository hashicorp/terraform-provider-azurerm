// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourceids

type ResourceId interface {

	// FromParseResult populates the Resource ID using the ParseResult provided in `input`
	FromParseResult(input ParseResult) error

	// ID returns the fully formatted ID for this Resource ID
	ID() string

	// String returns a friendly description of the components of this Resource ID
	// which is suitable for use in error messages (for example 'MyThing %q / Resource Group %q')
	String() string

	// Segments returns an ordered list of expected Segments that make up this Resource ID
	Segments() []Segment
}

type Segment struct {
	// ExampleValue is an example of a value for this field, which is intended only to
	// be used as a placeholder.
	ExampleValue string

	// FixedValue is the Fixed/Static value for this segment - only present when `Type`
	// is `ResourceProviderSegmentType` or `StaticSegmentType`.
	FixedValue *string

	// Name is the camelCased name of this segment, which is normalized (and safe to use as a
	// parameter/a field if necessary).
	Name string

	// PossibleValues is a list of possible values for this segment - only present when
	// `Type` is `ConstantSegmentType`
	PossibleValues *[]string

	// Type specifies the Type of Segment that this is, for example a `StaticSegmentType`
	Type SegmentType
}

type SegmentType string

const (
	// ConstantSegmentType specifies that this Segment is a Constant
	ConstantSegmentType SegmentType = "Constant"

	// ResourceGroupSegmentType specifies that this Segment is a Resource Group name
	ResourceGroupSegmentType SegmentType = "ResourceGroup"

	// ResourceProviderSegmentType specifies that this Segment is a Resource Provider
	ResourceProviderSegmentType SegmentType = "ResourceProvider"

	// ScopeSegmentType specifies that this Segment is a Scope
	ScopeSegmentType SegmentType = "Scope"

	// StaticSegmentType specifies that this Segment is a Static/Fixed Value
	StaticSegmentType SegmentType = "Static"

	// SubscriptionIdSegmentType specifies that this Segment is a Subscription ID
	SubscriptionIdSegmentType SegmentType = "SubscriptionId"

	// UserSpecifiedSegmentType specifies that this Segment is User-Specifiable
	UserSpecifiedSegmentType SegmentType = "UserSpecified"
)

// ConstantSegment is a helper which returns a Segment for a Constant
func ConstantSegment(name string, possibleValues []string, exampleValue string) Segment {
	return Segment{
		Name:           name,
		PossibleValues: &possibleValues,
		Type:           ConstantSegmentType,
		ExampleValue:   exampleValue,
	}
}

// ResourceGroupSegment is a helper which returns a Segment for a Resource Group
func ResourceGroupSegment(name, exampleValue string) Segment {
	return Segment{
		Name:         name,
		Type:         ResourceGroupSegmentType,
		ExampleValue: exampleValue,
	}
}

// ResourceProviderSegment is a helper which returns a Segment for a Resource Provider
func ResourceProviderSegment(name, resourceProvider, exampleValue string) Segment {
	return Segment{
		Name:         name,
		FixedValue:   &resourceProvider,
		Type:         ResourceProviderSegmentType,
		ExampleValue: exampleValue,
	}
}

// ScopeSegment is a helper which returns a Segment for a Scope
func ScopeSegment(name, exampleValue string) Segment {
	return Segment{
		Name:         name,
		Type:         ScopeSegmentType,
		ExampleValue: exampleValue,
	}
}

// StaticSegment is a helper which returns a Segment for a Static Value
func StaticSegment(name, staticValue, exampleValue string) Segment {
	return Segment{
		Name:         name,
		FixedValue:   &staticValue,
		Type:         StaticSegmentType,
		ExampleValue: exampleValue,
	}
}

// SubscriptionIdSegment is a helper which returns a Segment for a Subscription Id
func SubscriptionIdSegment(name, exampleValue string) Segment {
	return Segment{
		Name:         name,
		Type:         SubscriptionIdSegmentType,
		ExampleValue: exampleValue,
	}
}

// UserSpecifiedSegment is a helper which returns a Segment for a User Specified Segment
func UserSpecifiedSegment(name, exampleValue string) Segment {
	return Segment{
		Name:         name,
		Type:         UserSpecifiedSegmentType,
		ExampleValue: exampleValue,
	}
}
