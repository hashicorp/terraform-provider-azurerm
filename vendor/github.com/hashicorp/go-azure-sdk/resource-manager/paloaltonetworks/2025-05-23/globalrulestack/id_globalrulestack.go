package globalrulestack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&GlobalRulestackId{})
}

var _ resourceids.ResourceId = &GlobalRulestackId{}

// GlobalRulestackId is a struct representing the Resource ID for a Global Rulestack
type GlobalRulestackId struct {
	GlobalRulestackName string
}

// NewGlobalRulestackID returns a new GlobalRulestackId struct
func NewGlobalRulestackID(globalRulestackName string) GlobalRulestackId {
	return GlobalRulestackId{
		GlobalRulestackName: globalRulestackName,
	}
}

// ParseGlobalRulestackID parses 'input' into a GlobalRulestackId
func ParseGlobalRulestackID(input string) (*GlobalRulestackId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GlobalRulestackId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GlobalRulestackId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGlobalRulestackIDInsensitively parses 'input' case-insensitively into a GlobalRulestackId
// note: this method should only be used for API response data and not user input
func ParseGlobalRulestackIDInsensitively(input string) (*GlobalRulestackId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GlobalRulestackId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GlobalRulestackId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GlobalRulestackId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.GlobalRulestackName, ok = input.Parsed["globalRulestackName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "globalRulestackName", input)
	}

	return nil
}

// ValidateGlobalRulestackID checks that 'input' can be parsed as a Global Rulestack ID
func ValidateGlobalRulestackID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGlobalRulestackID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Global Rulestack ID
func (id GlobalRulestackId) ID() string {
	fmtString := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/%s"
	return fmt.Sprintf(fmtString, id.GlobalRulestackName)
}

// Segments returns a slice of Resource ID Segments which comprise this Global Rulestack ID
func (id GlobalRulestackId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticPaloAltoNetworksCloudngfw", "PaloAltoNetworks.Cloudngfw", "PaloAltoNetworks.Cloudngfw"),
		resourceids.StaticSegment("staticGlobalRulestacks", "globalRulestacks", "globalRulestacks"),
		resourceids.UserSpecifiedSegment("globalRulestackName", "globalRulestackName"),
	}
}

// String returns a human-readable description of this Global Rulestack ID
func (id GlobalRulestackId) String() string {
	components := []string{
		fmt.Sprintf("Global Rulestack Name: %q", id.GlobalRulestackName),
	}
	return fmt.Sprintf("Global Rulestack (%s)", strings.Join(components, "\n"))
}
