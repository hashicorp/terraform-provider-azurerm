package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CdnEndpointId struct {
	ResourceGroup string
	ProfileName   string
	Name          string
}

func NewCdnEndpointID(id CdnProfileId, name string) CdnEndpointId {
	return CdnEndpointId{
		ResourceGroup: id.ResourceGroup,
		ProfileName:   id.Name,
		Name:          name,
	}
}

func CdnEndpointID(input string) (*CdnEndpointId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse CDN Endpoint ID %q: %+v", input, err)
	}

	endpoint := CdnEndpointId{
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

func (id CdnEndpointId) ID(subscriptionId string) string {
	base := NewCdnProfileID(id.ResourceGroup, id.ProfileName).ID(subscriptionId)
	return fmt.Sprintf("%s/endpoints/%s", base, id.Name)
}
