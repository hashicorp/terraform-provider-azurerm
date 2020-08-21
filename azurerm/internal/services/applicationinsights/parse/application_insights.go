package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApplicationInsightsId struct {
	ResourceGroup string
	Name          string
}

func ApplicationInsightsID(input string) (*ApplicationInsightsId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Application Insights ID %q: %+v", input, err)
	}

	appId := ApplicationInsightsId{
		ResourceGroup: id.ResourceGroup,
	}

	if appId.Name, err = id.PopSegment("components"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &appId, nil
}

type ApplicationInsightsWebTestId struct {
	ResourceGroup string
	Name          string
}

func ApplicationInsightsWebTestID(input string) (*ApplicationInsightsWebTestId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Application Insights Web Test ID %q: %+v", input, err)
	}

	testid := ApplicationInsightsWebTestId{
		ResourceGroup: id.ResourceGroup,
	}

	if testid.Name, err = id.PopSegment("webtests"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &testid, nil
}
