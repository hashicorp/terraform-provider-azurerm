package volumequotarules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VolumeQuotaRuleId{})
}

var _ resourceids.ResourceId = &VolumeQuotaRuleId{}

// VolumeQuotaRuleId is a struct representing the Resource ID for a Volume Quota Rule
type VolumeQuotaRuleId struct {
	SubscriptionId      string
	ResourceGroupName   string
	NetAppAccountName   string
	CapacityPoolName    string
	VolumeName          string
	VolumeQuotaRuleName string
}

// NewVolumeQuotaRuleID returns a new VolumeQuotaRuleId struct
func NewVolumeQuotaRuleID(subscriptionId string, resourceGroupName string, netAppAccountName string, capacityPoolName string, volumeName string, volumeQuotaRuleName string) VolumeQuotaRuleId {
	return VolumeQuotaRuleId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		NetAppAccountName:   netAppAccountName,
		CapacityPoolName:    capacityPoolName,
		VolumeName:          volumeName,
		VolumeQuotaRuleName: volumeQuotaRuleName,
	}
}

// ParseVolumeQuotaRuleID parses 'input' into a VolumeQuotaRuleId
func ParseVolumeQuotaRuleID(input string) (*VolumeQuotaRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VolumeQuotaRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VolumeQuotaRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVolumeQuotaRuleIDInsensitively parses 'input' case-insensitively into a VolumeQuotaRuleId
// note: this method should only be used for API response data and not user input
func ParseVolumeQuotaRuleIDInsensitively(input string) (*VolumeQuotaRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VolumeQuotaRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VolumeQuotaRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VolumeQuotaRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetAppAccountName, ok = input.Parsed["netAppAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "netAppAccountName", input)
	}

	if id.CapacityPoolName, ok = input.Parsed["capacityPoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "capacityPoolName", input)
	}

	if id.VolumeName, ok = input.Parsed["volumeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "volumeName", input)
	}

	if id.VolumeQuotaRuleName, ok = input.Parsed["volumeQuotaRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "volumeQuotaRuleName", input)
	}

	return nil
}

// ValidateVolumeQuotaRuleID checks that 'input' can be parsed as a Volume Quota Rule ID
func ValidateVolumeQuotaRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVolumeQuotaRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Volume Quota Rule ID
func (id VolumeQuotaRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NetApp/netAppAccounts/%s/capacityPools/%s/volumes/%s/volumeQuotaRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName, id.VolumeQuotaRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Volume Quota Rule ID
func (id VolumeQuotaRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetApp", "Microsoft.NetApp", "Microsoft.NetApp"),
		resourceids.StaticSegment("staticNetAppAccounts", "netAppAccounts", "netAppAccounts"),
		resourceids.UserSpecifiedSegment("netAppAccountName", "netAppAccountName"),
		resourceids.StaticSegment("staticCapacityPools", "capacityPools", "capacityPools"),
		resourceids.UserSpecifiedSegment("capacityPoolName", "capacityPoolName"),
		resourceids.StaticSegment("staticVolumes", "volumes", "volumes"),
		resourceids.UserSpecifiedSegment("volumeName", "volumeName"),
		resourceids.StaticSegment("staticVolumeQuotaRules", "volumeQuotaRules", "volumeQuotaRules"),
		resourceids.UserSpecifiedSegment("volumeQuotaRuleName", "volumeQuotaRuleName"),
	}
}

// String returns a human-readable description of this Volume Quota Rule ID
func (id VolumeQuotaRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Net App Account Name: %q", id.NetAppAccountName),
		fmt.Sprintf("Capacity Pool Name: %q", id.CapacityPoolName),
		fmt.Sprintf("Volume Name: %q", id.VolumeName),
		fmt.Sprintf("Volume Quota Rule Name: %q", id.VolumeQuotaRuleName),
	}
	return fmt.Sprintf("Volume Quota Rule (%s)", strings.Join(components, "\n"))
}
