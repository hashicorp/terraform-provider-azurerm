package fqdnlistglobalrulestack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = FqdnListId{}

// FqdnListId is a struct representing the Resource ID for a Fqdn List
type FqdnListId struct {
	GlobalRulestackName string
	FqdnListName        string
}

// NewFqdnListID returns a new FqdnListId struct
func NewFqdnListID(globalRulestackName string, fqdnListName string) FqdnListId {
	return FqdnListId{
		GlobalRulestackName: globalRulestackName,
		FqdnListName:        fqdnListName,
	}
}

// ParseFqdnListID parses 'input' into a FqdnListId
func ParseFqdnListID(input string) (*FqdnListId, error) {
	parser := resourceids.NewParserFromResourceIdType(FqdnListId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FqdnListId{}

	if id.GlobalRulestackName, ok = parsed.Parsed["globalRulestackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "globalRulestackName", *parsed)
	}

	if id.FqdnListName, ok = parsed.Parsed["fqdnListName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fqdnListName", *parsed)
	}

	return &id, nil
}

// ParseFqdnListIDInsensitively parses 'input' case-insensitively into a FqdnListId
// note: this method should only be used for API response data and not user input
func ParseFqdnListIDInsensitively(input string) (*FqdnListId, error) {
	parser := resourceids.NewParserFromResourceIdType(FqdnListId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FqdnListId{}

	if id.GlobalRulestackName, ok = parsed.Parsed["globalRulestackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "globalRulestackName", *parsed)
	}

	if id.FqdnListName, ok = parsed.Parsed["fqdnListName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fqdnListName", *parsed)
	}

	return &id, nil
}

// ValidateFqdnListID checks that 'input' can be parsed as a Fqdn List ID
func ValidateFqdnListID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFqdnListID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Fqdn List ID
func (id FqdnListId) ID() string {
	fmtString := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/%s/fqdnLists/%s"
	return fmt.Sprintf(fmtString, id.GlobalRulestackName, id.FqdnListName)
}

// Segments returns a slice of Resource ID Segments which comprise this Fqdn List ID
func (id FqdnListId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticPaloAltoNetworksCloudngfw", "PaloAltoNetworks.Cloudngfw", "PaloAltoNetworks.Cloudngfw"),
		resourceids.StaticSegment("staticGlobalRulestacks", "globalRulestacks", "globalRulestacks"),
		resourceids.UserSpecifiedSegment("globalRulestackName", "globalRulestackValue"),
		resourceids.StaticSegment("staticFqdnLists", "fqdnLists", "fqdnLists"),
		resourceids.UserSpecifiedSegment("fqdnListName", "fqdnListValue"),
	}
}

// String returns a human-readable description of this Fqdn List ID
func (id FqdnListId) String() string {
	components := []string{
		fmt.Sprintf("Global Rulestack Name: %q", id.GlobalRulestackName),
		fmt.Sprintf("Fqdn List Name: %q", id.FqdnListName),
	}
	return fmt.Sprintf("Fqdn List (%s)", strings.Join(components, "\n"))
}
