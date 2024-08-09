package agentversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&OsTypeId{})
}

var _ resourceids.ResourceId = &OsTypeId{}

// OsTypeId is a struct representing the Resource ID for a Os Type
type OsTypeId struct {
	OsTypeName string
}

// NewOsTypeID returns a new OsTypeId struct
func NewOsTypeID(osTypeName string) OsTypeId {
	return OsTypeId{
		OsTypeName: osTypeName,
	}
}

// ParseOsTypeID parses 'input' into a OsTypeId
func ParseOsTypeID(input string) (*OsTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OsTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OsTypeId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOsTypeIDInsensitively parses 'input' case-insensitively into a OsTypeId
// note: this method should only be used for API response data and not user input
func ParseOsTypeIDInsensitively(input string) (*OsTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OsTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OsTypeId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OsTypeId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.OsTypeName, ok = input.Parsed["osTypeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "osTypeName", input)
	}

	return nil
}

// ValidateOsTypeID checks that 'input' can be parsed as a Os Type ID
func ValidateOsTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOsTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Os Type ID
func (id OsTypeId) ID() string {
	fmtString := "/providers/Microsoft.HybridCompute/osType/%s"
	return fmt.Sprintf(fmtString, id.OsTypeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Os Type ID
func (id OsTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHybridCompute", "Microsoft.HybridCompute", "Microsoft.HybridCompute"),
		resourceids.StaticSegment("staticOsType", "osType", "osType"),
		resourceids.UserSpecifiedSegment("osTypeName", "osTypeValue"),
	}
}

// String returns a human-readable description of this Os Type ID
func (id OsTypeId) String() string {
	components := []string{
		fmt.Sprintf("Os Type Name: %q", id.OsTypeName),
	}
	return fmt.Sprintf("Os Type (%s)", strings.Join(components, "\n"))
}
