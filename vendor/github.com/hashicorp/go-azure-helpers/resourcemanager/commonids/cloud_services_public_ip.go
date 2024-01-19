// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &CloudServicesPublicIPAddressId{}

// CloudServicesPublicIPAddressId is a struct representing the Resource ID for a Cloud Services Public I P Address
type CloudServicesPublicIPAddressId struct {
	SubscriptionId       string
	ResourceGroupName    string
	CloudServiceName     string
	RoleInstanceName     string
	NetworkInterfaceName string
	IpConfigurationName  string
	PublicIPAddressName  string
}

// NewCloudServicesPublicIPAddressID returns a new CloudServicesPublicIPAddressId struct
func NewCloudServicesPublicIPAddressID(subscriptionId string, resourceGroupName string, cloudServiceName string, roleInstanceName string, networkInterfaceName string, ipConfigurationName string, publicIPAddressName string) CloudServicesPublicIPAddressId {
	return CloudServicesPublicIPAddressId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		CloudServiceName:     cloudServiceName,
		RoleInstanceName:     roleInstanceName,
		NetworkInterfaceName: networkInterfaceName,
		IpConfigurationName:  ipConfigurationName,
		PublicIPAddressName:  publicIPAddressName,
	}
}

// ParseCloudServicesPublicIPAddressID parses 'input' into a CloudServicesPublicIPAddressId
func ParseCloudServicesPublicIPAddressID(input string) (*CloudServicesPublicIPAddressId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CloudServicesPublicIPAddressId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CloudServicesPublicIPAddressId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCloudServicesPublicIPAddressIDInsensitively parses 'input' case-insensitively into a CloudServicesPublicIPAddressId
// note: this method should only be used for API response data and not user input
func ParseCloudServicesPublicIPAddressIDInsensitively(input string) (*CloudServicesPublicIPAddressId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CloudServicesPublicIPAddressId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CloudServicesPublicIPAddressId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CloudServicesPublicIPAddressId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CloudServiceName, ok = input.Parsed["cloudServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "cloudServiceName", input)
	}

	if id.RoleInstanceName, ok = input.Parsed["roleInstanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleInstanceName", input)
	}

	if id.NetworkInterfaceName, ok = input.Parsed["networkInterfaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkInterfaceName", input)
	}

	if id.IpConfigurationName, ok = input.Parsed["ipConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "ipConfigurationName", input)
	}

	if id.PublicIPAddressName, ok = input.Parsed["publicIPAddressName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "publicIPAddressName", input)
	}

	return nil
}

// ValidateCloudServicesPublicIPAddressID checks that 'input' can be parsed as a Cloud Services Public I P Address ID
func ValidateCloudServicesPublicIPAddressID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCloudServicesPublicIPAddressID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cloud Services Public IP Address ID
func (id CloudServicesPublicIPAddressId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/cloudServices/%s/roleInstances/%s/networkInterfaces/%s/ipConfigurations/%s/publicIPAddresses/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CloudServiceName, id.RoleInstanceName, id.NetworkInterfaceName, id.IpConfigurationName, id.PublicIPAddressName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cloud Services Public I P Address ID
func (id CloudServicesPublicIPAddressId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("cloudServices", "cloudServices", "cloudServices"),
		resourceids.UserSpecifiedSegment("cloudServiceName", "cloudServiceValue"),
		resourceids.StaticSegment("roleInstances", "roleInstances", "roleInstances"),
		resourceids.UserSpecifiedSegment("roleInstanceName", "roleInstanceValue"),
		resourceids.StaticSegment("networkInterfaces", "networkInterfaces", "networkInterfaces"),
		resourceids.UserSpecifiedSegment("networkInterfaceName", "networkInterfaceValue"),
		resourceids.StaticSegment("ipConfigurations", "ipConfigurations", "ipConfigurations"),
		resourceids.UserSpecifiedSegment("ipConfigurationName", "ipConfigurationValue"),
		resourceids.StaticSegment("publicIPAddresses", "publicIPAddresses", "publicIPAddresses"),
		resourceids.UserSpecifiedSegment("publicIPAddressName", "publicIPAddressValue"),
	}
}

// String returns a human-readable description of this Cloud Services Public I P Address ID
func (id CloudServicesPublicIPAddressId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cloud Service Name: %q", id.CloudServiceName),
		fmt.Sprintf("Role Instance Name: %q", id.RoleInstanceName),
		fmt.Sprintf("Network Interface Name: %q", id.NetworkInterfaceName),
		fmt.Sprintf("Ip Configuration Name: %q", id.IpConfigurationName),
		fmt.Sprintf("Public I P Address Name: %q", id.PublicIPAddressName),
	}
	return fmt.Sprintf("Cloud Services Public IP Address (%s)", strings.Join(components, "\n"))
}
