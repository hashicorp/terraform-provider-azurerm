package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CdnEndpointId struct {
	SubscriptionId string
	ResourceGroup  string
	ProfileName    string
	Name           string
}

func NewCdnEndpointID(id CdnProfileId, name string) CdnEndpointId {
	return CdnEndpointId{
		SubscriptionId: id.SubscriptionId,
		ResourceGroup:  id.ResourceGroup,
		ProfileName:    id.Name,
		Name:           name,
	}
}

func (id CdnEndpointId) ID(_ string) string {
	base := NewCdnProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName)
	return fmt.Sprintf("%s/endpoints/%s", base, id.Name)
}

func CdnEndpointID(input string) (*CdnEndpointId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse CDN Endpoint ID %q: %+v", input, err)
	}

	endpoint := CdnEndpointId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
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
