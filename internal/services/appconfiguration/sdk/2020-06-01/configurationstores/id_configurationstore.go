package configurationstores

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ConfigurationStoreId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewConfigurationStoreID(subscriptionId, resourceGroup, name string) ConfigurationStoreId {
	return ConfigurationStoreId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id ConfigurationStoreId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Configuration Store", segmentsStr)
}

func (id ConfigurationStoreId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppConfiguration/configurationStores/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// ConfigurationStoreID parses a ConfigurationStore ID into an ConfigurationStoreId struct
func ConfigurationStoreID(input string) (*ConfigurationStoreId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ConfigurationStoreId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("configurationStores"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ConfigurationStoreIDInsensitively parses an ConfigurationStore ID into an ConfigurationStoreId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ConfigurationStoreID method should be used instead for validation etc.
func ConfigurationStoreIDInsensitively(input string) (*ConfigurationStoreId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ConfigurationStoreId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'configurationStores' segment
	configurationStoresKey := "configurationStores"
	for key := range id.Path {
		if strings.EqualFold(key, configurationStoresKey) {
			configurationStoresKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(configurationStoresKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
