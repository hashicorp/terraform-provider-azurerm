package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SecurityCenterSettingId struct {
	SettingName string
}

func SecurityCenterSettingID(input string) (*SecurityCenterSettingId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse Security Center Setting ID %q: %+v", input, err)
	}

	setting := SecurityCenterSettingId{}

	if setting.SettingName, err = id.PopSegment("settings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &setting, nil
}
