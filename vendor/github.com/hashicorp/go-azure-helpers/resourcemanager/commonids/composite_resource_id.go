// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// CompositeResourceID is a struct representing the Resource ID for a Composite Resource Id
type CompositeResourceID[T1 resourceids.ResourceId, T2 resourceids.ResourceId] struct {
	// First specifies the first component of this Resource ID.
	// This is in the format `{first}|{second}`.
	First T1

	// Second specifies the second component of this Resource ID
	// This is in the format `{first}|{second}`.
	Second T2
}

// NewCompositeResourceID returns a new CompositeResourceID struct
func NewCompositeResourceID[T1 resourceids.ResourceId, T2 resourceids.ResourceId](first T1, second T2) CompositeResourceID[T1, T2] {
	return CompositeResourceID[T1, T2]{
		First:  first,
		Second: second,
	}
}

// ID returns the formatted Composite Resource Id
func (id CompositeResourceID[T1, T2]) ID() string {
	fmtString := "%s|%s"
	return fmt.Sprintf(fmtString, id.First.ID(), id.Second.ID())
}

// String returns a human-readable description of this Composite Resource Id
func (id CompositeResourceID[T1, T2]) String() string {
	fmtString := "Composite Resource ID (%s | %s)"
	return fmt.Sprintf(fmtString, id.First.String(), id.Second.String())
}

// ParseCompositeResourceID parses 'input' and two ResourceIds (first,second) into a CompositeResourceID
// The 'input' should be a string containing 2 resource ids separated by "|"
// eg: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group|/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Sql/servers/serverValue"
// The first and second ResourceIds should match the types in the 'input' string in the order in which they appear
// eg:
//
//	input := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group|/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.Sql/servers/serverValue"
//	first := ResourceGroupId{}
//	second := SqlServerId{}
//	id, err := ParseCompositeResourceID(input, &first, &second)
func ParseCompositeResourceID[T1 resourceids.ResourceId, T2 resourceids.ResourceId](input string, first T1, second T2) (*CompositeResourceID[T1, T2], error) {
	return parseCompositeResourceID(input, first, second, false)
}

// ParseCompositeResourceIDInsensitively parses 'input' and two ResourceIds (first,second) case-insensitively into a CompositeResourceID
// note: this method should only be used for API response data and not user input
func ParseCompositeResourceIDInsensitively[T1 resourceids.ResourceId, T2 resourceids.ResourceId](input string, first T1, second T2) (*CompositeResourceID[T1, T2], error) {
	return parseCompositeResourceID(input, first, second, true)
}

func parseCompositeResourceID[T1 resourceids.ResourceId, T2 resourceids.ResourceId](input string, first T1, second T2, insensitively bool) (*CompositeResourceID[T1, T2], error) {

	components := strings.Split(input, "|")
	if len(components) != 2 {
		return nil, fmt.Errorf("expected 2 resourceids but got %d", len(components))
	}

	output := CompositeResourceID[T1, T2]{
		First:  first,
		Second: second,
	}

	// Parse the first of the two Resource IDs from the components
	firstParser := resourceids.NewParserFromResourceIdType(output.First)
	firstParseResult, err := firstParser.Parse(components[0], insensitively)
	if err != nil {
		return nil, fmt.Errorf("parsing first id part %q of CompositeResourceID: %v", components[0], err)
	}
	err = output.First.FromParseResult(*firstParseResult)
	if err != nil {
		return nil, fmt.Errorf("populating first id part %q of CompositeResourceID: %v", components[0], err)
	}

	// Parse the second of the two Resource IDs from the components
	secondParser := resourceids.NewParserFromResourceIdType(output.Second)
	secondParseResult, err := secondParser.Parse(components[1], insensitively)
	if err != nil {
		return nil, fmt.Errorf("parsing second id part %q of CompositeResourceID: %v", components[1], err)
	}
	err = output.Second.FromParseResult(*secondParseResult)
	if err != nil {
		return nil, fmt.Errorf("populating second id part %q of CompositeResourceID: %v", components[1], err)
	}

	return &output, nil
}
