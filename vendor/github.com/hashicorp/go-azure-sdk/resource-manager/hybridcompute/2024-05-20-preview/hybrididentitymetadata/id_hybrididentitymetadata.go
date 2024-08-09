package hybrididentitymetadata

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&HybridIdentityMetadataId{})
}

var _ resourceids.ResourceId = &HybridIdentityMetadataId{}

// HybridIdentityMetadataId is a struct representing the Resource ID for a Hybrid Identity Metadata
type HybridIdentityMetadataId struct {
	SubscriptionId             string
	ResourceGroupName          string
	MachineName                string
	HybridIdentityMetadataName string
}

// NewHybridIdentityMetadataID returns a new HybridIdentityMetadataId struct
func NewHybridIdentityMetadataID(subscriptionId string, resourceGroupName string, machineName string, hybridIdentityMetadataName string) HybridIdentityMetadataId {
	return HybridIdentityMetadataId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		MachineName:                machineName,
		HybridIdentityMetadataName: hybridIdentityMetadataName,
	}
}

// ParseHybridIdentityMetadataID parses 'input' into a HybridIdentityMetadataId
func ParseHybridIdentityMetadataID(input string) (*HybridIdentityMetadataId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HybridIdentityMetadataId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HybridIdentityMetadataId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHybridIdentityMetadataIDInsensitively parses 'input' case-insensitively into a HybridIdentityMetadataId
// note: this method should only be used for API response data and not user input
func ParseHybridIdentityMetadataIDInsensitively(input string) (*HybridIdentityMetadataId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HybridIdentityMetadataId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HybridIdentityMetadataId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HybridIdentityMetadataId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MachineName, ok = input.Parsed["machineName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "machineName", input)
	}

	if id.HybridIdentityMetadataName, ok = input.Parsed["hybridIdentityMetadataName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hybridIdentityMetadataName", input)
	}

	return nil
}

// ValidateHybridIdentityMetadataID checks that 'input' can be parsed as a Hybrid Identity Metadata ID
func ValidateHybridIdentityMetadataID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHybridIdentityMetadataID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Hybrid Identity Metadata ID
func (id HybridIdentityMetadataId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HybridCompute/machines/%s/hybridIdentityMetadata/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MachineName, id.HybridIdentityMetadataName)
}

// Segments returns a slice of Resource ID Segments which comprise this Hybrid Identity Metadata ID
func (id HybridIdentityMetadataId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHybridCompute", "Microsoft.HybridCompute", "Microsoft.HybridCompute"),
		resourceids.StaticSegment("staticMachines", "machines", "machines"),
		resourceids.UserSpecifiedSegment("machineName", "machineValue"),
		resourceids.StaticSegment("staticHybridIdentityMetadata", "hybridIdentityMetadata", "hybridIdentityMetadata"),
		resourceids.UserSpecifiedSegment("hybridIdentityMetadataName", "hybridIdentityMetadataValue"),
	}
}

// String returns a human-readable description of this Hybrid Identity Metadata ID
func (id HybridIdentityMetadataId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Machine Name: %q", id.MachineName),
		fmt.Sprintf("Hybrid Identity Metadata Name: %q", id.HybridIdentityMetadataName),
	}
	return fmt.Sprintf("Hybrid Identity Metadata (%s)", strings.Join(components, "\n"))
}
