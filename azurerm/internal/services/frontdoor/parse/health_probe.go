package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type HealthProbeId struct {
	SubscriptionId         string
	ResourceGroup          string
	FrontDoorName          string
	HealthProbeSettingName string
}

func NewHealthProbeID(subscriptionId, resourceGroup, frontDoorName, healthProbeSettingName string) HealthProbeId {
	return HealthProbeId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		FrontDoorName:          frontDoorName,
		HealthProbeSettingName: healthProbeSettingName,
	}
}

func (id HealthProbeId) String() string {
	segments := []string{
		fmt.Sprintf("Health Probe Setting Name %q", id.HealthProbeSettingName),
		fmt.Sprintf("Front Door Name %q", id.FrontDoorName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Health Probe", segmentsStr)
}

func (id HealthProbeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontDoors/%s/healthProbeSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FrontDoorName, id.HealthProbeSettingName)
}

// HealthProbeID parses a HealthProbe ID into an HealthProbeId struct
func HealthProbeID(input string) (*HealthProbeId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := HealthProbeId{
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
	if resourceId.HealthProbeSettingName, err = id.PopSegment("healthProbeSettings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// HealthProbeIDInsensitively parses an HealthProbe ID into an HealthProbeId struct, insensitively
// This should only be used to parse an ID for rewriting, the HealthProbeID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func HealthProbeIDInsensitively(input string) (*HealthProbeId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := HealthProbeId{
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

	// find the correct casing for the 'healthProbeSettings' segment
	healthProbeSettingsKey := "healthProbeSettings"
	for key := range id.Path {
		if strings.EqualFold(key, healthProbeSettingsKey) {
			healthProbeSettingsKey = key
			break
		}
	}
	if resourceId.HealthProbeSettingName, err = id.PopSegment(healthProbeSettingsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
