package projectconnectionresource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProjectConnectionId{})
}

var _ resourceids.ResourceId = &ProjectConnectionId{}

// ProjectConnectionId is a struct representing the Resource ID for a Project Connection
type ProjectConnectionId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	ProjectName       string
	ConnectionName    string
}

// NewProjectConnectionID returns a new ProjectConnectionId struct
func NewProjectConnectionID(subscriptionId string, resourceGroupName string, accountName string, projectName string, connectionName string) ProjectConnectionId {
	return ProjectConnectionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		ProjectName:       projectName,
		ConnectionName:    connectionName,
	}
}

// ParseProjectConnectionID parses 'input' into a ProjectConnectionId
func ParseProjectConnectionID(input string) (*ProjectConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProjectConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProjectConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProjectConnectionIDInsensitively parses 'input' case-insensitively into a ProjectConnectionId
// note: this method should only be used for API response data and not user input
func ParseProjectConnectionIDInsensitively(input string) (*ProjectConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProjectConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProjectConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProjectConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AccountName, ok = input.Parsed["accountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "accountName", input)
	}

	if id.ProjectName, ok = input.Parsed["projectName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "projectName", input)
	}

	if id.ConnectionName, ok = input.Parsed["connectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "connectionName", input)
	}

	return nil
}

// ValidateProjectConnectionID checks that 'input' can be parsed as a Project Connection ID
func ValidateProjectConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProjectConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Project Connection ID
func (id ProjectConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CognitiveServices/accounts/%s/projects/%s/connections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ProjectName, id.ConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Project Connection ID
func (id ProjectConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCognitiveServices", "Microsoft.CognitiveServices", "Microsoft.CognitiveServices"),
		resourceids.StaticSegment("staticAccounts", "accounts", "accounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountName"),
		resourceids.StaticSegment("staticProjects", "projects", "projects"),
		resourceids.UserSpecifiedSegment("projectName", "projectName"),
		resourceids.StaticSegment("staticConnections", "connections", "connections"),
		resourceids.UserSpecifiedSegment("connectionName", "connectionName"),
	}
}

// String returns a human-readable description of this Project Connection ID
func (id ProjectConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Project Name: %q", id.ProjectName),
		fmt.Sprintf("Connection Name: %q", id.ConnectionName),
	}
	return fmt.Sprintf("Project Connection (%s)", strings.Join(components, "\n"))
}
