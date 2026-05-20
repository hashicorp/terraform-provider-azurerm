package attacheddatabaseconfigurations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AttachedDatabaseConfigurationId{})
}

var _ resourceids.ResourceId = &AttachedDatabaseConfigurationId{}

// AttachedDatabaseConfigurationId is a struct representing the Resource ID for a Attached Database Configuration
type AttachedDatabaseConfigurationId struct {
	SubscriptionId                    string
	ResourceGroupName                 string
	ClusterName                       string
	AttachedDatabaseConfigurationName string
}

// NewAttachedDatabaseConfigurationID returns a new AttachedDatabaseConfigurationId struct
func NewAttachedDatabaseConfigurationID(subscriptionId string, resourceGroupName string, clusterName string, attachedDatabaseConfigurationName string) AttachedDatabaseConfigurationId {
	return AttachedDatabaseConfigurationId{
		SubscriptionId:                    subscriptionId,
		ResourceGroupName:                 resourceGroupName,
		ClusterName:                       clusterName,
		AttachedDatabaseConfigurationName: attachedDatabaseConfigurationName,
	}
}

// ParseAttachedDatabaseConfigurationID parses 'input' into a AttachedDatabaseConfigurationId
func ParseAttachedDatabaseConfigurationID(input string) (*AttachedDatabaseConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AttachedDatabaseConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AttachedDatabaseConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAttachedDatabaseConfigurationIDInsensitively parses 'input' case-insensitively into a AttachedDatabaseConfigurationId
// note: this method should only be used for API response data and not user input
func ParseAttachedDatabaseConfigurationIDInsensitively(input string) (*AttachedDatabaseConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AttachedDatabaseConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AttachedDatabaseConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AttachedDatabaseConfigurationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ClusterName, ok = input.Parsed["clusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "clusterName", input)
	}

	if id.AttachedDatabaseConfigurationName, ok = input.Parsed["attachedDatabaseConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "attachedDatabaseConfigurationName", input)
	}

	return nil
}

// ValidateAttachedDatabaseConfigurationID checks that 'input' can be parsed as a Attached Database Configuration ID
func ValidateAttachedDatabaseConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAttachedDatabaseConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Attached Database Configuration ID
func (id AttachedDatabaseConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/clusters/%s/attachedDatabaseConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.AttachedDatabaseConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Attached Database Configuration ID
func (id AttachedDatabaseConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKusto", "Microsoft.Kusto", "Microsoft.Kusto"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterName"),
		resourceids.StaticSegment("staticAttachedDatabaseConfigurations", "attachedDatabaseConfigurations", "attachedDatabaseConfigurations"),
		resourceids.UserSpecifiedSegment("attachedDatabaseConfigurationName", "attachedDatabaseConfigurationName"),
	}
}

// String returns a human-readable description of this Attached Database Configuration ID
func (id AttachedDatabaseConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Attached Database Configuration Name: %q", id.AttachedDatabaseConfigurationName),
	}
	return fmt.Sprintf("Attached Database Configuration (%s)", strings.Join(components, "\n"))
}
