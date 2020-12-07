package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LoadBalancingId struct {
	SubscriptionId           string
	ResourceGroup            string
	FrontDoorName            string
	LoadBalancingSettingName string
}

func NewLoadBalancingID(subscriptionId, resourceGroup, frontDoorName, loadBalancingSettingName string) LoadBalancingId {
	return LoadBalancingId{
		SubscriptionId:           subscriptionId,
		ResourceGroup:            resourceGroup,
		FrontDoorName:            frontDoorName,
		LoadBalancingSettingName: loadBalancingSettingName,
	}
}

func (id LoadBalancingId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Front Door Name %q", id.FrontDoorName),
		fmt.Sprintf("Load Balancing Setting Name %q", id.LoadBalancingSettingName),
	}
	return strings.Join(segments, " / ")
}

func (id LoadBalancingId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontDoors/%s/loadBalancingSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FrontDoorName, id.LoadBalancingSettingName)
}

// LoadBalancingID parses a LoadBalancing ID into an LoadBalancingId struct
func LoadBalancingID(input string) (*LoadBalancingId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LoadBalancingId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FrontDoorName, err = id.PopSegment("frontDoors"); err != nil {
		return nil, err
	}
	if resourceId.LoadBalancingSettingName, err = id.PopSegment("loadBalancingSettings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// LoadBalancingIDInsensitively parses an LoadBalancing ID into an LoadBalancingId struct, insensitively
// This should only be used to parse an ID for rewriting, the LoadBalancingID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func LoadBalancingIDInsensitively(input string) (*LoadBalancingId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LoadBalancingId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'frontDoors' segment
	frontDoorsKey := "frontDoors"
	for key := range id.Path {
		if strings.EqualFold(key, frontDoorsKey) {
			frontDoorsKey = key
			break
		}
	}
	if resourceId.FrontDoorName, err = id.PopSegment(frontDoorsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'loadBalancingSettings' segment
	loadBalancingSettingsKey := "loadBalancingSettings"
	for key := range id.Path {
		if strings.EqualFold(key, loadBalancingSettingsKey) {
			loadBalancingSettingsKey = key
			break
		}
	}
	if resourceId.LoadBalancingSettingName, err = id.PopSegment(loadBalancingSettingsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
