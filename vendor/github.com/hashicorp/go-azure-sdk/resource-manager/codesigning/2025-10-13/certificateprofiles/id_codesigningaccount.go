package certificateprofiles

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CodeSigningAccountId{})
}

var _ resourceids.ResourceId = &CodeSigningAccountId{}

// CodeSigningAccountId is a struct representing the Resource ID for a Code Signing Account
type CodeSigningAccountId struct {
	SubscriptionId         string
	ResourceGroupName      string
	CodeSigningAccountName string
}

// NewCodeSigningAccountID returns a new CodeSigningAccountId struct
func NewCodeSigningAccountID(subscriptionId string, resourceGroupName string, codeSigningAccountName string) CodeSigningAccountId {
	return CodeSigningAccountId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		CodeSigningAccountName: codeSigningAccountName,
	}
}

// ParseCodeSigningAccountID parses 'input' into a CodeSigningAccountId
func ParseCodeSigningAccountID(input string) (*CodeSigningAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CodeSigningAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CodeSigningAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCodeSigningAccountIDInsensitively parses 'input' case-insensitively into a CodeSigningAccountId
// note: this method should only be used for API response data and not user input
func ParseCodeSigningAccountIDInsensitively(input string) (*CodeSigningAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CodeSigningAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CodeSigningAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CodeSigningAccountId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CodeSigningAccountName, ok = input.Parsed["codeSigningAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "codeSigningAccountName", input)
	}

	return nil
}

// ValidateCodeSigningAccountID checks that 'input' can be parsed as a Code Signing Account ID
func ValidateCodeSigningAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCodeSigningAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Code Signing Account ID
func (id CodeSigningAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CodeSigning/codeSigningAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CodeSigningAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Code Signing Account ID
func (id CodeSigningAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCodeSigning", "Microsoft.CodeSigning", "Microsoft.CodeSigning"),
		resourceids.StaticSegment("staticCodeSigningAccounts", "codeSigningAccounts", "codeSigningAccounts"),
		resourceids.UserSpecifiedSegment("codeSigningAccountName", "codeSigningAccountName"),
	}
}

// String returns a human-readable description of this Code Signing Account ID
func (id CodeSigningAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Code Signing Account Name: %q", id.CodeSigningAccountName),
	}
	return fmt.Sprintf("Code Signing Account (%s)", strings.Join(components, "\n"))
}
