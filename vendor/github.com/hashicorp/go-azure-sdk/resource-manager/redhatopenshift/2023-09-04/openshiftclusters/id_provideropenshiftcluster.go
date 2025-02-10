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
	recaser.RegisterResourceId(&ProviderOpenShiftClusterId{})
}

var _ resourceids.ResourceId = &ProviderOpenShiftClusterId{}

// ProviderOpenShiftClusterId is a struct representing the Resource ID for a Provider Open Shift Cluster
type ProviderOpenShiftClusterId struct {
	SubscriptionId       string
	ResourceGroupName    string
	OpenShiftClusterName string
}

// NewProviderOpenShiftClusterID returns a new ProviderOpenShiftClusterId struct
func NewProviderOpenShiftClusterID(subscriptionId string, resourceGroupName string, openShiftClusterName string) ProviderOpenShiftClusterId {
	return ProviderOpenShiftClusterId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		OpenShiftClusterName: openShiftClusterName,
	}
}

// ParseProviderOpenShiftClusterID parses 'input' into a ProviderOpenShiftClusterId
func ParseProviderOpenShiftClusterID(input string) (*ProviderOpenShiftClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderOpenShiftClusterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderOpenShiftClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProviderOpenShiftClusterIDInsensitively parses 'input' case-insensitively into a ProviderOpenShiftClusterId
// note: this method should only be used for API response data and not user input
func ParseProviderOpenShiftClusterIDInsensitively(input string) (*ProviderOpenShiftClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderOpenShiftClusterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderOpenShiftClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProviderOpenShiftClusterId) FromParseResult(input resourceids.ParseResult) error {
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

// ValidateProviderOpenShiftClusterID checks that 'input' can be parsed as a Provider Open Shift Cluster ID
func ValidateProviderOpenShiftClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderOpenShiftClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Open Shift Cluster ID
func (id ProviderOpenShiftClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RedHatOpenShift/openShiftClusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.OpenShiftClusterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Open Shift Cluster ID
func (id ProviderOpenShiftClusterId) Segments() []resourceids.Segment {
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

// String returns a human-readable description of this Provider Open Shift Cluster ID
func (id ProviderOpenShiftClusterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Open Shift Cluster Name: %q", id.OpenShiftClusterName),
	}
	return fmt.Sprintf("Provider Open Shift Cluster (%s)", strings.Join(components, "\n"))
}
