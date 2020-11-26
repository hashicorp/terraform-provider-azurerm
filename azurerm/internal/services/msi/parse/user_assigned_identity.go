package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type UserAssignedIdentityId struct {
	// "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/contoso-resource-group/providers/Microsoft.ManagedIdentity/userAssignedIdentities/contoso-identity"
	Subscription  string
	ResourceGroup string
	Name          string
}

func UserAssignedIdentityID(input string) (*UserAssignedIdentityId, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("Bad: UserAssignedIdentityId cannot be an empty string")
	}

	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	userAssignedIdentityId := UserAssignedIdentityId{
		Subscription:  id.SubscriptionID,
		ResourceGroup: id.ResourceGroup,
	}

	if name, err := id.PopSegment("userAssignedIdentities"); err != nil {
		return nil, fmt.Errorf("Bad: missing userAssignedIdentities segment in ID (%q)", input)
	} else {
		userAssignedIdentityId.Name = name
	}

	return &userAssignedIdentityId, nil
}
