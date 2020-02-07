package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AdvisorResGroupId struct {
	ResourceGroup string
}

func AdvisorSubscriptionID(input string) error {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse Advisor subscription ID %q: %+v", input, err)
	}

	if id.ResourceGroup != "" {
		return fmt.Errorf("[ERROR] There should be no resource group in Advisor subscription ID %q", input)
	}

	if _, err = id.PopSegment("configurations"); err != nil {
		return err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return err
	}

	return nil
}

func AdvisorResGroupID(input string) (*AdvisorResGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Advisor ResGroup ID %q: %+v", input, err)
	}

	advisorResGroup := AdvisorResGroupId{
		ResourceGroup: id.ResourceGroup,
	}

	if _, err = id.PopSegment("configurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &advisorResGroup, nil
}
