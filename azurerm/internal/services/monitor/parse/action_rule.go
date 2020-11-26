package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ActionRuleId struct {
	ResourceGroup string
	Name          string
}

func ActionRuleID(input string) (*ActionRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Action Rule ID %q: %+v", input, err)
	}

	actionRule := ActionRuleId{
		ResourceGroup: id.ResourceGroup,
	}
	if actionRule.Name, err = id.PopSegment("actionRules"); err != nil {
		return nil, fmt.Errorf("parsing Action Rule ID %q: %+v", input, err)
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, fmt.Errorf("parsing ActionRule ID %q: %+v", input, err)
	}

	return &actionRule, nil
}
