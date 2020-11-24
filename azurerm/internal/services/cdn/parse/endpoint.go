package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EndpointId struct {
	ResourceGroup string
	ProfileName   string
	EndpointName  string
}

func NewEndpointID(resourceGroup, profileName, name string) EndpointId {
	return EndpointId{
		ResourceGroup: resourceGroup,
		ProfileName:   profileName,
		EndpointName:  name,
	}
}

func EndpointID(input string) (*EndpointId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse CDN Endpoint ID %q: %+v", input, err)
	}

	endpoint := EndpointId{
		ResourceGroup: id.ResourceGroup,
	}

	if endpoint.ProfileName, err = id.PopSegment("profiles"); err != nil {
		return nil, err
	}

	if endpoint.EndpointName, err = id.PopSegment("endpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &endpoint, nil
}

func (id EndpointId) ID(subscriptionId string) string {
	base := NewProfileID(id.ResourceGroup, id.ProfileName).ID(subscriptionId)
	return fmt.Sprintf("%s/endpoints/%s", base, id.EndpointName)
}
