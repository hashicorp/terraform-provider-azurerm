package apidefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DefinitionId{})
}

var _ resourceids.ResourceId = &DefinitionId{}

// DefinitionId is a struct representing the Resource ID for a Definition
type DefinitionId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceName     string
	ApiName           string
	VersionName       string
	DefinitionName    string
}

// NewDefinitionID returns a new DefinitionId struct
func NewDefinitionID(subscriptionId string, resourceGroupName string, serviceName string, workspaceName string, apiName string, versionName string, definitionName string) DefinitionId {
	return DefinitionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceName:     workspaceName,
		ApiName:           apiName,
		VersionName:       versionName,
		DefinitionName:    definitionName,
	}
}

// ParseDefinitionID parses 'input' into a DefinitionId
func ParseDefinitionID(input string) (*DefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDefinitionIDInsensitively parses 'input' case-insensitively into a DefinitionId
// note: this method should only be used for API response data and not user input
func ParseDefinitionIDInsensitively(input string) (*DefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DefinitionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if id.WorkspaceName, ok = input.Parsed["workspaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", input)
	}

	if id.ApiName, ok = input.Parsed["apiName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "apiName", input)
	}

	if id.VersionName, ok = input.Parsed["versionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "versionName", input)
	}

	if id.DefinitionName, ok = input.Parsed["definitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "definitionName", input)
	}

	return nil
}

// ValidateDefinitionID checks that 'input' can be parsed as a Definition ID
func ValidateDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Definition ID
func (id DefinitionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiCenter/services/%s/workspaces/%s/apis/%s/versions/%s/definitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceName, id.ApiName, id.VersionName, id.DefinitionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Definition ID
func (id DefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiCenter", "Microsoft.ApiCenter", "Microsoft.ApiCenter"),
		resourceids.StaticSegment("staticServices", "services", "services"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceName"),
		resourceids.StaticSegment("staticApis", "apis", "apis"),
		resourceids.UserSpecifiedSegment("apiName", "apiName"),
		resourceids.StaticSegment("staticVersions", "versions", "versions"),
		resourceids.UserSpecifiedSegment("versionName", "versionName"),
		resourceids.StaticSegment("staticDefinitions", "definitions", "definitions"),
		resourceids.UserSpecifiedSegment("definitionName", "definitionName"),
	}
}

// String returns a human-readable description of this Definition ID
func (id DefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Api Name: %q", id.ApiName),
		fmt.Sprintf("Version Name: %q", id.VersionName),
		fmt.Sprintf("Definition Name: %q", id.DefinitionName),
	}
	return fmt.Sprintf("Definition (%s)", strings.Join(components, "\n"))
}
