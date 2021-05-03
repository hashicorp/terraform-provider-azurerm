package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CustomHttpsConfigurationId struct {
	SubscriptionId               string
	ResourceGroup                string
	FrontDoorName                string
	CustomHttpsConfigurationName string
}

func NewCustomHttpsConfigurationID(subscriptionId, resourceGroup, frontDoorName, customHttpsConfigurationName string) CustomHttpsConfigurationId {
	return CustomHttpsConfigurationId{
		SubscriptionId:               subscriptionId,
		ResourceGroup:                resourceGroup,
		FrontDoorName:                frontDoorName,
		CustomHttpsConfigurationName: customHttpsConfigurationName,
	}
}

func (id CustomHttpsConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Custom Https Configuration Name %q", id.CustomHttpsConfigurationName),
		fmt.Sprintf("Front Door Name %q", id.FrontDoorName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Custom Https Configuration", segmentsStr)
}

func (id CustomHttpsConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontDoors/%s/customHttpsConfiguration/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FrontDoorName, id.CustomHttpsConfigurationName)
}

// CustomHttpsConfigurationID parses a CustomHttpsConfiguration ID into an CustomHttpsConfigurationId struct
func CustomHttpsConfigurationID(input string) (*CustomHttpsConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CustomHttpsConfigurationId{
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
	if resourceId.CustomHttpsConfigurationName, err = id.PopSegment("customHttpsConfiguration"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// CustomHttpsConfigurationIDInsensitively parses an CustomHttpsConfiguration ID into an CustomHttpsConfigurationId struct, insensitively
// This should only be used to parse an ID for rewriting, the CustomHttpsConfigurationID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func CustomHttpsConfigurationIDInsensitively(input string) (*CustomHttpsConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CustomHttpsConfigurationId{
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

	// find the correct casing for the 'customHttpsConfiguration' segment
	customHttpsConfigurationKey := "customHttpsConfiguration"
	for key := range id.Path {
		if strings.EqualFold(key, customHttpsConfigurationKey) {
			customHttpsConfigurationKey = key
			break
		}
	}
	if resourceId.CustomHttpsConfigurationName, err = id.PopSegment(customHttpsConfigurationKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
