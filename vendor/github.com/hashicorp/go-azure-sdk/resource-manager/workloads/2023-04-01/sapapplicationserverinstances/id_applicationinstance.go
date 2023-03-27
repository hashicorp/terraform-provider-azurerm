package sapapplicationserverinstances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ApplicationInstanceId{}

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
	parser := resourceids.NewParserFromResourceIdType(ApplicationInstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApplicationInstanceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.SapVirtualInstanceName, ok = parsed.Parsed["sapVirtualInstanceName"]; !ok {
		return nil, fmt.Errorf("the segment 'sapVirtualInstanceName' was not found in the resource id %q", input)
	}

	if id.ApplicationInstanceName, ok = parsed.Parsed["applicationInstanceName"]; !ok {
		return nil, fmt.Errorf("the segment 'applicationInstanceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseApplicationInstanceIDInsensitively parses 'input' case-insensitively into a ApplicationInstanceId
// note: this method should only be used for API response data and not user input
func ParseApplicationInstanceIDInsensitively(input string) (*ApplicationInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApplicationInstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApplicationInstanceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.SapVirtualInstanceName, ok = parsed.Parsed["sapVirtualInstanceName"]; !ok {
		return nil, fmt.Errorf("the segment 'sapVirtualInstanceName' was not found in the resource id %q", input)
	}

	if id.ApplicationInstanceName, ok = parsed.Parsed["applicationInstanceName"]; !ok {
		return nil, fmt.Errorf("the segment 'applicationInstanceName' was not found in the resource id %q", input)
	}

	return &id, nil
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
		resourceids.UserSpecifiedSegment("sapVirtualInstanceName", "sapVirtualInstanceValue"),
		resourceids.StaticSegment("staticApplicationInstances", "applicationInstances", "applicationInstances"),
		resourceids.UserSpecifiedSegment("applicationInstanceName", "applicationInstanceValue"),
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
