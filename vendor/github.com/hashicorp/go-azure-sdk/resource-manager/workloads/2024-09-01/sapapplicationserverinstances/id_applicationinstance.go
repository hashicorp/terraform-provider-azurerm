package sapapplicationserverinstances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ApplicationInstanceId{})
}

var _ resourceids.ResourceId = &ApplicationInstanceId{}

// ApplicationInstanceId is a struct representing the Resource ID for a Application Instance
type ApplicationInstanceId struct {
	SubscriptionId          string
	ResourceGroupName       string
	SapVirtualInstanceName  string
	ApplicationInstanceName string
}

// NewApplicationInstanceID returns a new ApplicationInstanceId struct
func NewApplicationInstanceID(subscriptionId string, resourceGroupName string, sapVirtualInstanceName string, applicationInstanceName string) ApplicationInstanceId {
	return ApplicationInstanceId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		SapVirtualInstanceName:  sapVirtualInstanceName,
		ApplicationInstanceName: applicationInstanceName,
	}
}

// ParseApplicationInstanceID parses 'input' into a ApplicationInstanceId
func ParseApplicationInstanceID(input string) (*ApplicationInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationInstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApplicationInstanceIDInsensitively parses 'input' case-insensitively into a ApplicationInstanceId
// note: this method should only be used for API response data and not user input
func ParseApplicationInstanceIDInsensitively(input string) (*ApplicationInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationInstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApplicationInstanceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SapVirtualInstanceName, ok = input.Parsed["sapVirtualInstanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "sapVirtualInstanceName", input)
	}

	if id.ApplicationInstanceName, ok = input.Parsed["applicationInstanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "applicationInstanceName", input)
	}

	return nil
}

// ValidateApplicationInstanceID checks that 'input' can be parsed as a Application Instance ID
func ValidateApplicationInstanceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApplicationInstanceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Application Instance ID
func (id ApplicationInstanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Workloads/sapVirtualInstances/%s/applicationInstances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SapVirtualInstanceName, id.ApplicationInstanceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Application Instance ID
func (id ApplicationInstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWorkloads", "Microsoft.Workloads", "Microsoft.Workloads"),
		resourceids.StaticSegment("staticSapVirtualInstances", "sapVirtualInstances", "sapVirtualInstances"),
		resourceids.UserSpecifiedSegment("sapVirtualInstanceName", "sapVirtualInstanceName"),
		resourceids.StaticSegment("staticApplicationInstances", "applicationInstances", "applicationInstances"),
		resourceids.UserSpecifiedSegment("applicationInstanceName", "applicationInstanceName"),
	}
}

// String returns a human-readable description of this Application Instance ID
func (id ApplicationInstanceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Sap Virtual Instance Name: %q", id.SapVirtualInstanceName),
		fmt.Sprintf("Application Instance Name: %q", id.ApplicationInstanceName),
	}
	return fmt.Sprintf("Application Instance (%s)", strings.Join(components, "\n"))
}
