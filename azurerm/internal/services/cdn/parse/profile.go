package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ProfileId struct {
	ResourceGroup string
	ProfileName   string
}

func NewProfileID(resourceGroup, name string) ProfileId {
	return ProfileId{
		ResourceGroup: resourceGroup,
		ProfileName:   name,
	}
}

func ProfileID(input string) (*ProfileId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse CDN Profile ID %q: %+v", input, err)
	}

	profile := ProfileId{
		ResourceGroup: id.ResourceGroup,
	}

	if profile.ProfileName, err = id.PopSegment("profiles"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (id ProfileId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s",
		subscriptionId, id.ResourceGroup, id.ProfileName)
}
