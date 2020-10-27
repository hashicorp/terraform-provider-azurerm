package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CdnProfileId struct {
	ResourceGroup string
	Name          string
}

func NewCdnProfileID(resourceGroup, name string) CdnProfileId {
	return CdnProfileId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func CdnProfileID(input string) (*CdnProfileId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse CDN Profile ID %q: %+v", input, err)
	}

	profile := CdnProfileId{
		ResourceGroup: id.ResourceGroup,
	}

	if profile.Name, err = id.PopSegment("profiles"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (id CdnProfileId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s",
		subscriptionId, id.ResourceGroup, id.Name)
}
