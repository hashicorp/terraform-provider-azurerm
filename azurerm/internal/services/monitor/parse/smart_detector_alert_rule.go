package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SmartDetectorAlertRuleId struct {
	ResourceGroup string
	Name          string
}

func SmartDetectorAlertRuleID(input string) (*SmartDetectorAlertRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Smart Detector Alert Rule ID %q: %+v", input, err)
	}

	smartDetectorAlertRuleId := SmartDetectorAlertRuleId{
		ResourceGroup: id.ResourceGroup,
	}
	if smartDetectorAlertRuleId.Name, err = id.PopSegment("smartdetectoralertrules"); err != nil {
		return nil, fmt.Errorf("parsing Smart Detector Alert Rule ID %q: %+v", input, err)
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, fmt.Errorf("parsing Smart Detector Alert Rule ID %q: %+v", input, err)
	}

	return &smartDetectorAlertRuleId, nil
}
