package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SubnetId struct {
	SubscriptionId     string
	ResourceGroup      string
	VirtualNetworkName string
	Name               string
}

func NewSubnetID(subscriptionId, resourceGroup, virtualNetworkName, name string) SubnetId {
	return SubnetId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		VirtualNetworkName: virtualNetworkName,
		Name:               name,
	}
}

func (id SubnetId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Virtual Network Name %q", id.VirtualNetworkName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Subnet", segmentsStr)
}

func (id SubnetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualNetworkName, id.Name)
}

// SubnetID parses a Subnet ID into an SubnetId struct
func SubnetID(input string) (*SubnetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SubnetId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VirtualNetworkName, err = id.PopSegment("virtualNetworks"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("subnets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SubnetIDInsensitively parses an Subnet ID into an SubnetId struct, insensitively
// This should only be used to parse an ID for rewriting, the SubnetID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SubnetIDInsensitively(input string) (*SubnetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SubnetId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'virtualNetworks' segment
	virtualNetworksKey := "virtualNetworks"
	for key := range id.Path {
		if strings.EqualFold(key, virtualNetworksKey) {
			virtualNetworksKey = key
			break
		}
	}
	if resourceId.VirtualNetworkName, err = id.PopSegment(virtualNetworksKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'subnets' segment
	subnetsKey := "subnets"
	for key := range id.Path {
		if strings.EqualFold(key, subnetsKey) {
			subnetsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(subnetsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
