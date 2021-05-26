package configurationstores

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type PrivateLinkResourceId struct {
	SubscriptionId         string
	ResourceGroup          string
	ConfigurationStoreName string
	Name                   string
}

func NewPrivateLinkResourceID(subscriptionId, resourceGroup, configurationStoreName, name string) PrivateLinkResourceId {
	return PrivateLinkResourceId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		ConfigurationStoreName: configurationStoreName,
		Name:                   name,
	}
}

func (id PrivateLinkResourceId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Configuration Store Name %q", id.ConfigurationStoreName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Private Link Resource", segmentsStr)
}

func (id PrivateLinkResourceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppConfiguration/configurationStores/%s/privateLinkResources/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ConfigurationStoreName, id.Name)
}

// PrivateLinkResourceID parses a PrivateLinkResource ID into an PrivateLinkResourceId struct
func PrivateLinkResourceID(input string) (*PrivateLinkResourceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PrivateLinkResourceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ConfigurationStoreName, err = id.PopSegment("configurationStores"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("privateLinkResources"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// PrivateLinkResourceIDInsensitively parses an PrivateLinkResource ID into an PrivateLinkResourceId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the PrivateLinkResourceID method should be used instead for validation etc.
func PrivateLinkResourceIDInsensitively(input string) (*PrivateLinkResourceId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PrivateLinkResourceId{
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
	if resourceId.ConfigurationStoreName, err = id.PopSegment(configurationStoresKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'privateLinkResources' segment
	privateLinkResourcesKey := "privateLinkResources"
	for key := range id.Path {
		if strings.EqualFold(key, privateLinkResourcesKey) {
			privateLinkResourcesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(privateLinkResourcesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
