package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CdnEndpointId struct {
	Subscription  string
	Provider      string
	ResourceGroup string
	ProfileName   string
	Name          string
}

func CdnEndpointID(input string) (*CdnEndpointId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse CDN Endpoint ID %q: %+v", input, err)
	}

	endpoint := CdnEndpointId{
		Subscription:  id.SubscriptionID,
		Provider:      id.Provider,
		ResourceGroup: id.ResourceGroup,
	}

	if endpoint.ProfileName, err = id.PopSegment("profiles"); err != nil {
		return nil, err
	}

	if endpoint.Name, err = id.PopSegment("endpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &endpoint, nil
}

// This ID is a workaround for issue: https://github.com/Azure/azure-rest-api-specs/issues/10576
func (id CdnEndpointId) ID() string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/%s/profiles/%s/endpoints/%s",
		id.Subscription, id.ResourceGroup, id.Provider, id.ProfileName, id.Name)
}
