package sapcentralinstances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CentralInstanceId{})
}

var _ resourceids.ResourceId = &CentralInstanceId{}

// CentralInstanceId is a struct representing the Resource ID for a Central Instance
type CentralInstanceId struct {
	SubscriptionId         string
	ResourceGroupName      string
	SapVirtualInstanceName string
	CentralInstanceName    string
}

// NewCentralInstanceID returns a new CentralInstanceId struct
func NewCentralInstanceID(subscriptionId string, resourceGroupName string, sapVirtualInstanceName string, centralInstanceName string) CentralInstanceId {
	return CentralInstanceId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		SapVirtualInstanceName: sapVirtualInstanceName,
		CentralInstanceName:    centralInstanceName,
	}
}

// ParseCentralInstanceID parses 'input' into a CentralInstanceId
func ParseCentralInstanceID(input string) (*CentralInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CentralInstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CentralInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCentralInstanceIDInsensitively parses 'input' case-insensitively into a CentralInstanceId
// note: this method should only be used for API response data and not user input
func ParseCentralInstanceIDInsensitively(input string) (*CentralInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CentralInstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CentralInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CentralInstanceId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.CentralInstanceName, ok = input.Parsed["centralInstanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "centralInstanceName", input)
	}

	return nil
}

// ValidateCentralInstanceID checks that 'input' can be parsed as a Central Instance ID
func ValidateCentralInstanceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCentralInstanceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Central Instance ID
func (id CentralInstanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Workloads/sapVirtualInstances/%s/centralInstances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SapVirtualInstanceName, id.CentralInstanceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Central Instance ID
func (id CentralInstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWorkloads", "Microsoft.Workloads", "Microsoft.Workloads"),
		resourceids.StaticSegment("staticSapVirtualInstances", "sapVirtualInstances", "sapVirtualInstances"),
		resourceids.UserSpecifiedSegment("sapVirtualInstanceName", "sapVirtualInstanceName"),
		resourceids.StaticSegment("staticCentralInstances", "centralInstances", "centralInstances"),
		resourceids.UserSpecifiedSegment("centralInstanceName", "centralInstanceName"),
	}
}

// String returns a human-readable description of this Central Instance ID
func (id CentralInstanceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Sap Virtual Instance Name: %q", id.SapVirtualInstanceName),
		fmt.Sprintf("Central Instance Name: %q", id.CentralInstanceName),
	}
	return fmt.Sprintf("Central Instance (%s)", strings.Join(components, "\n"))
}
