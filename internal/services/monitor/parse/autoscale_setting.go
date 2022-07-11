package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AutoscaleSettingId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewAutoscaleSettingID(subscriptionId, resourceGroup, name string) AutoscaleSettingId {
	return AutoscaleSettingId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id AutoscaleSettingId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Autoscale Setting", segmentsStr)
}

func (id AutoscaleSettingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/autoscaleSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// AutoscaleSettingID parses a AutoscaleSetting ID into an AutoscaleSettingId struct
func AutoscaleSettingID(input string) (*AutoscaleSettingId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AutoscaleSettingId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("autoscaleSettings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// AutoscaleSettingIDInsensitively parses an AutoscaleSetting ID into an AutoscaleSettingId struct, insensitively
// This should only be used to parse an ID for rewriting, the AutoscaleSettingID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func AutoscaleSettingIDInsensitively(input string) (*AutoscaleSettingId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AutoscaleSettingId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'autoscaleSettings' segment
	autoscaleSettingsKey := "autoscaleSettings"
	for key := range id.Path {
		if strings.EqualFold(key, autoscaleSettingsKey) {
			autoscaleSettingsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(autoscaleSettingsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
