package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SecurityCenterSettingId struct {
	SubscriptionId string
	SettingName    string
}

func NewSecurityCenterSettingID(subscriptionId, settingName string) SecurityCenterSettingId {
	return SecurityCenterSettingId{
		SubscriptionId: subscriptionId,
		SettingName:    settingName,
	}
}

func (id SecurityCenterSettingId) String() string {
	segments := []string{
		fmt.Sprintf("Setting Name %q", id.SettingName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Security Center Setting", segmentsStr)
}

func (id SecurityCenterSettingId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Security/settings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.SettingName)
}

// SecurityCenterSettingID parses a SecurityCenterSetting ID into an SecurityCenterSettingId struct
func SecurityCenterSettingID(input string) (*SecurityCenterSettingId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SecurityCenterSettingId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.SettingName, err = id.PopSegment("settings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
