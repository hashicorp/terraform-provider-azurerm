package openshiftclusters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&OpenShiftClusterId{})
}

var _ resourceids.ResourceId = &OpenShiftClusterId{}

// OpenShiftClusterId is a struct representing the Resource ID for a Open Shift Cluster
type OpenShiftClusterId struct {
	SubscriptionId       string
	ResourceGroupName    string
	OpenShiftClusterName string
}

// NewOpenShiftClusterID returns a new OpenShiftClusterId struct
func NewOpenShiftClusterID(subscriptionId string, resourceGroupName string, openShiftClusterName string) OpenShiftClusterId {
	return OpenShiftClusterId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		OpenShiftClusterName: openShiftClusterName,
	}
}

// ParseOpenShiftClusterID parses 'input' into a OpenShiftClusterId
func ParseOpenShiftClusterID(input string) (*OpenShiftClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OpenShiftClusterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OpenShiftClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOpenShiftClusterIDInsensitively parses 'input' case-insensitively into a OpenShiftClusterId
// note: this method should only be used for API response data and not user input
func ParseOpenShiftClusterIDInsensitively(input string) (*OpenShiftClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OpenShiftClusterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OpenShiftClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OpenShiftClusterId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.OpenShiftClusterName, ok = input.Parsed["openShiftClusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "openShiftClusterName", input)
	}

	return nil
}

// ValidateOpenShiftClusterID checks that 'input' can be parsed as a Open Shift Cluster ID
func ValidateOpenShiftClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOpenShiftClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Open Shift Cluster ID
func (id OpenShiftClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RedHatOpenShift/openShiftClusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.OpenShiftClusterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Open Shift Cluster ID
func (id OpenShiftClusterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftRedHatOpenShift", "Microsoft.RedHatOpenShift", "Microsoft.RedHatOpenShift"),
		resourceids.StaticSegment("staticOpenShiftClusters", "openShiftClusters", "openShiftClusters"),
		resourceids.UserSpecifiedSegment("openShiftClusterName", "openShiftClusterName"),
	}
}

// String returns a human-readable description of this Open Shift Cluster ID
func (id OpenShiftClusterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Open Shift Cluster Name: %q", id.OpenShiftClusterName),
	}
	return fmt.Sprintf("Open Shift Cluster (%s)", strings.Join(components, "\n"))
}
