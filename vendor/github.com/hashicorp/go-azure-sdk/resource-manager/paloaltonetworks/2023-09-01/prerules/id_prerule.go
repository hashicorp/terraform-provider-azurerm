package prerules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PreRuleId{})
}

var _ resourceids.ResourceId = &PreRuleId{}

// PreRuleId is a struct representing the Resource ID for a Pre Rule
type PreRuleId struct {
	GlobalRulestackName string
	PreRuleName         string
}

// NewPreRuleID returns a new PreRuleId struct
func NewPreRuleID(globalRulestackName string, preRuleName string) PreRuleId {
	return PreRuleId{
		GlobalRulestackName: globalRulestackName,
		PreRuleName:         preRuleName,
	}
}

// ParsePreRuleID parses 'input' into a PreRuleId
func ParsePreRuleID(input string) (*PreRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PreRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PreRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePreRuleIDInsensitively parses 'input' case-insensitively into a PreRuleId
// note: this method should only be used for API response data and not user input
func ParsePreRuleIDInsensitively(input string) (*PreRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PreRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PreRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PreRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.GlobalRulestackName, ok = input.Parsed["globalRulestackName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "globalRulestackName", input)
	}

	if id.PreRuleName, ok = input.Parsed["preRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "preRuleName", input)
	}

	return nil
}

// ValidatePreRuleID checks that 'input' can be parsed as a Pre Rule ID
func ValidatePreRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePreRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Pre Rule ID
func (id PreRuleId) ID() string {
	fmtString := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/%s/preRules/%s"
	return fmt.Sprintf(fmtString, id.GlobalRulestackName, id.PreRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Pre Rule ID
func (id PreRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticPaloAltoNetworksCloudngfw", "PaloAltoNetworks.Cloudngfw", "PaloAltoNetworks.Cloudngfw"),
		resourceids.StaticSegment("staticGlobalRulestacks", "globalRulestacks", "globalRulestacks"),
		resourceids.UserSpecifiedSegment("globalRulestackName", "globalRulestackName"),
		resourceids.StaticSegment("staticPreRules", "preRules", "preRules"),
		resourceids.UserSpecifiedSegment("preRuleName", "preRuleName"),
	}
}

// String returns a human-readable description of this Pre Rule ID
func (id PreRuleId) String() string {
	components := []string{
		fmt.Sprintf("Global Rulestack Name: %q", id.GlobalRulestackName),
		fmt.Sprintf("Pre Rule Name: %q", id.PreRuleName),
	}
	return fmt.Sprintf("Pre Rule (%s)", strings.Join(components, "\n"))
}
