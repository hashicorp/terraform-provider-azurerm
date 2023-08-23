package postrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PostRuleId{}

// PostRuleId is a struct representing the Resource ID for a Post Rule
type PostRuleId struct {
	GlobalRulestackName string
	PostRuleName        string
}

// NewPostRuleID returns a new PostRuleId struct
func NewPostRuleID(globalRulestackName string, postRuleName string) PostRuleId {
	return PostRuleId{
		GlobalRulestackName: globalRulestackName,
		PostRuleName:        postRuleName,
	}
}

// ParsePostRuleID parses 'input' into a PostRuleId
func ParsePostRuleID(input string) (*PostRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(PostRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PostRuleId{}

	if id.GlobalRulestackName, ok = parsed.Parsed["globalRulestackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "globalRulestackName", *parsed)
	}

	if id.PostRuleName, ok = parsed.Parsed["postRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "postRuleName", *parsed)
	}

	return &id, nil
}

// ParsePostRuleIDInsensitively parses 'input' case-insensitively into a PostRuleId
// note: this method should only be used for API response data and not user input
func ParsePostRuleIDInsensitively(input string) (*PostRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(PostRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PostRuleId{}

	if id.GlobalRulestackName, ok = parsed.Parsed["globalRulestackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "globalRulestackName", *parsed)
	}

	if id.PostRuleName, ok = parsed.Parsed["postRuleName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "postRuleName", *parsed)
	}

	return &id, nil
}

// ValidatePostRuleID checks that 'input' can be parsed as a Post Rule ID
func ValidatePostRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePostRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Post Rule ID
func (id PostRuleId) ID() string {
	fmtString := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/%s/postRules/%s"
	return fmt.Sprintf(fmtString, id.GlobalRulestackName, id.PostRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Post Rule ID
func (id PostRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticPaloAltoNetworksCloudngfw", "PaloAltoNetworks.Cloudngfw", "PaloAltoNetworks.Cloudngfw"),
		resourceids.StaticSegment("staticGlobalRulestacks", "globalRulestacks", "globalRulestacks"),
		resourceids.UserSpecifiedSegment("globalRulestackName", "globalRulestackValue"),
		resourceids.StaticSegment("staticPostRules", "postRules", "postRules"),
		resourceids.UserSpecifiedSegment("postRuleName", "postRuleValue"),
	}
}

// String returns a human-readable description of this Post Rule ID
func (id PostRuleId) String() string {
	components := []string{
		fmt.Sprintf("Global Rulestack Name: %q", id.GlobalRulestackName),
		fmt.Sprintf("Post Rule Name: %q", id.PostRuleName),
	}
	return fmt.Sprintf("Post Rule (%s)", strings.Join(components, "\n"))
}
