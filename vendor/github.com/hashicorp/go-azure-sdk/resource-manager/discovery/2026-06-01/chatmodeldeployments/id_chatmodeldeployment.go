package chatmodeldeployments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ChatModelDeploymentId{})
}

var _ resourceids.ResourceId = &ChatModelDeploymentId{}

// ChatModelDeploymentId is a struct representing the Resource ID for a Chat Model Deployment
type ChatModelDeploymentId struct {
	SubscriptionId          string
	ResourceGroupName       string
	WorkspaceName           string
	ChatModelDeploymentName string
}

// NewChatModelDeploymentID returns a new ChatModelDeploymentId struct
func NewChatModelDeploymentID(subscriptionId string, resourceGroupName string, workspaceName string, chatModelDeploymentName string) ChatModelDeploymentId {
	return ChatModelDeploymentId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		WorkspaceName:           workspaceName,
		ChatModelDeploymentName: chatModelDeploymentName,
	}
}

// ParseChatModelDeploymentID parses 'input' into a ChatModelDeploymentId
func ParseChatModelDeploymentID(input string) (*ChatModelDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ChatModelDeploymentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ChatModelDeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseChatModelDeploymentIDInsensitively parses 'input' case-insensitively into a ChatModelDeploymentId
// note: this method should only be used for API response data and not user input
func ParseChatModelDeploymentIDInsensitively(input string) (*ChatModelDeploymentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ChatModelDeploymentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ChatModelDeploymentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ChatModelDeploymentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.WorkspaceName, ok = input.Parsed["workspaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", input)
	}

	if id.ChatModelDeploymentName, ok = input.Parsed["chatModelDeploymentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "chatModelDeploymentName", input)
	}

	return nil
}

// ValidateChatModelDeploymentID checks that 'input' can be parsed as a Chat Model Deployment ID
func ValidateChatModelDeploymentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseChatModelDeploymentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Chat Model Deployment ID
func (id ChatModelDeploymentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Discovery/workspaces/%s/chatModelDeployments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.ChatModelDeploymentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Chat Model Deployment ID
func (id ChatModelDeploymentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDiscovery", "Microsoft.Discovery", "Microsoft.Discovery"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceName"),
		resourceids.StaticSegment("staticChatModelDeployments", "chatModelDeployments", "chatModelDeployments"),
		resourceids.UserSpecifiedSegment("chatModelDeploymentName", "chatModelDeploymentName"),
	}
}

// String returns a human-readable description of this Chat Model Deployment ID
func (id ChatModelDeploymentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Chat Model Deployment Name: %q", id.ChatModelDeploymentName),
	}
	return fmt.Sprintf("Chat Model Deployment (%s)", strings.Join(components, "\n"))
}
