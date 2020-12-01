package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApiManagementPolicyId struct {
	ResourceGroup string
	ServiceName   string
	Name          string
}

func ApiManagementPolicyID(input string) (*ApiManagementPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	policy := ApiManagementPolicyId{
		ResourceGroup: id.ResourceGroup,
	}

	if policy.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}

	if policy.Name, err = id.PopSegment("policies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &policy, nil
}
