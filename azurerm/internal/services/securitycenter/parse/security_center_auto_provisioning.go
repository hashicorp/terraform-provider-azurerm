package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SecurityCenterAutoProvisioningId struct {
	SubscriptionId              string
	AutoProvisioningSettingName string
}

func NewSecurityCenterAutoProvisioningID(subscriptionId, autoProvisioningSettingName string) SecurityCenterAutoProvisioningId {
	return SecurityCenterAutoProvisioningId{
		SubscriptionId:              subscriptionId,
		AutoProvisioningSettingName: autoProvisioningSettingName,
	}
}

func (id SecurityCenterAutoProvisioningId) String() string {
	segments := []string{
		fmt.Sprintf("Auto Provisioning Setting Name %q", id.AutoProvisioningSettingName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Security Center Auto Provisioning", segmentsStr)
}

func (id SecurityCenterAutoProvisioningId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Security/autoProvisioningSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.AutoProvisioningSettingName)
}

// SecurityCenterAutoProvisioningID parses a SecurityCenterAutoProvisioning ID into an SecurityCenterAutoProvisioningId struct
func SecurityCenterAutoProvisioningID(input string) (*SecurityCenterAutoProvisioningId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SecurityCenterAutoProvisioningId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.AutoProvisioningSettingName, err = id.PopSegment("autoProvisioningSettings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
