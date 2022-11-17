package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AutomanageConfigurationProfileId struct {
	SubscriptionId           string
	ResourceGroup            string
	ConfigurationProfileName string
}

func NewAutomanageConfigurationProfileID(subscriptionId, resourceGroup, configurationProfileName string) AutomanageConfigurationProfileId {
	return AutomanageConfigurationProfileId{
		SubscriptionId:           subscriptionId,
		ResourceGroup:            resourceGroup,
		ConfigurationProfileName: configurationProfileName,
	}
}

func (id AutomanageConfigurationProfileId) String() string {
	segments := []string{
		fmt.Sprintf("Configuration Profile Name %q", id.ConfigurationProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Automanage Configuration Profile", segmentsStr)
}

func (id AutomanageConfigurationProfileId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automanage/configurationProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ConfigurationProfileName)
}

// AutomanageConfigurationProfileID parses a AutomanageConfigurationProfile ID into an AutomanageConfigurationProfileId struct
func AutomanageConfigurationProfileID(input string) (*AutomanageConfigurationProfileId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AutomanageConfigurationProfileId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ConfigurationProfileName, err = id.PopSegment("configurationProfiles"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
