package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AdvancedThreatProtectionId struct {
	SubscriptionId                      string
	AdvancedThreatProtectionSettingName string
}

func NewAdvancedThreatProtectionID(subscriptionId, advancedThreatProtectionSettingName string) AdvancedThreatProtectionId {
	return AdvancedThreatProtectionId{
		SubscriptionId:                      subscriptionId,
		AdvancedThreatProtectionSettingName: advancedThreatProtectionSettingName,
	}
}

func (id AdvancedThreatProtectionId) String() string {
	segments := []string{
		fmt.Sprintf("Advanced Threat Protection Setting Name %q", id.AdvancedThreatProtectionSettingName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Advanced Threat Protection", segmentsStr)
}

func (id AdvancedThreatProtectionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Security/advancedThreatProtectionSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.AdvancedThreatProtectionSettingName)
}

// AdvancedThreatProtectionID parses a AdvancedThreatProtection ID into an AdvancedThreatProtectionId struct
func AdvancedThreatProtectionID(input string) (*AdvancedThreatProtectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AdvancedThreatProtectionId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.AdvancedThreatProtectionSettingName, err = id.PopSegment("advancedThreatProtectionSettings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
