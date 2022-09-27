package parse

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type AutomanageConfigurationProfileId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewAutomanageConfigurationProfileID(subscriptionId string, resourcegroup string, name string) AutomanageConfigurationProfileId {
	return AutomanageConfigurationProfileId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourcegroup,
		Name:           name,
	}
}

func (id AutomanageConfigurationProfileId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automanage/configurationProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func AutomanageConfigurationProfileID(input string) (*AutomanageConfigurationProfileId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing automanageConfigurationProfile ID %q: %+v", input, err)
	}

	automanageConfigurationProfile := AutomanageConfigurationProfileId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}
	if automanageConfigurationProfile.Name, err = id.PopSegment("configurationProfiles"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &automanageConfigurationProfile, nil
}
