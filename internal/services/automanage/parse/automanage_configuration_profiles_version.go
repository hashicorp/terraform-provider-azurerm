package parse

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type AutomanageConfigurationProfilesVersionId struct {
	SubscriptionId           string
	ResourceGroup            string
	ConfigurationProfileName string
	Name                     string
}

func NewAutomanageConfigurationProfilesVersionID(subscriptionId string, resourcegroup string, configurationprofilename string, name string) AutomanageConfigurationProfilesVersionId {
	return AutomanageConfigurationProfilesVersionId{
		SubscriptionId:           subscriptionId,
		ResourceGroup:            resourcegroup,
		ConfigurationProfileName: configurationprofilename,
		Name:                     name,
	}
}

func (id AutomanageConfigurationProfilesVersionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automanage/configurationProfiles/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ConfigurationProfileName, id.Name)
}

func AutomanageConfigurationProfilesVersionID(input string) (*AutomanageConfigurationProfilesVersionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing automanageConfigurationProfilesVersion ID %q: %+v", input, err)
	}

	automanageConfigurationProfilesVersion := AutomanageConfigurationProfilesVersionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}
	if automanageConfigurationProfilesVersion.ConfigurationProfileName, err = id.PopSegment("configurationProfiles"); err != nil {
		return nil, err
	}
	if automanageConfigurationProfilesVersion.Name, err = id.PopSegment("versions"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &automanageConfigurationProfilesVersion, nil
}
