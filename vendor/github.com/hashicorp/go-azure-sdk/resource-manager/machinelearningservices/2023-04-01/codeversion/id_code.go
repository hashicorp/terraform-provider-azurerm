package codeversion

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = CodeId{}

// CodeId is a struct representing the Resource ID for a Code
type CodeId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	CodeName          string
}

// NewCodeID returns a new CodeId struct
func NewCodeID(subscriptionId string, resourceGroupName string, workspaceName string, codeName string) CodeId {
	return CodeId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		CodeName:          codeName,
	}
}

// ParseCodeID parses 'input' into a CodeId
func ParseCodeID(input string) (*CodeId, error) {
	parser := resourceids.NewParserFromResourceIdType(CodeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CodeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.CodeName, ok = parsed.Parsed["codeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "codeName", *parsed)
	}

	return &id, nil
}

// ParseCodeIDInsensitively parses 'input' case-insensitively into a CodeId
// note: this method should only be used for API response data and not user input
func ParseCodeIDInsensitively(input string) (*CodeId, error) {
	parser := resourceids.NewParserFromResourceIdType(CodeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CodeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.CodeName, ok = parsed.Parsed["codeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "codeName", *parsed)
	}

	return &id, nil
}

// ValidateCodeID checks that 'input' can be parsed as a Code ID
func ValidateCodeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCodeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Code ID
func (id CodeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MachineLearningServices/workspaces/%s/codes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.CodeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Code ID
func (id CodeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMachineLearningServices", "Microsoft.MachineLearningServices", "Microsoft.MachineLearningServices"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticCodes", "codes", "codes"),
		resourceids.UserSpecifiedSegment("codeName", "codeValue"),
	}
}

// String returns a human-readable description of this Code ID
func (id CodeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Code Name: %q", id.CodeName),
	}
	return fmt.Sprintf("Code (%s)", strings.Join(components, "\n"))
}
