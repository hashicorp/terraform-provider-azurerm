package securitycenter

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AdvancedThreatProtectionId struct {
	TargetResourceID string
	SettingName      string
}

func NewAdvancedThreatProtectionId(targetResourceId string) AdvancedThreatProtectionId {
	return AdvancedThreatProtectionId{
		TargetResourceID: targetResourceId,
		SettingName:      "current",
	}
}

func (id AdvancedThreatProtectionId) ID(_ string) string {
	fmtString := "%s/providers/Microsoft.Security/advancedThreatProtectionSettings/%s"
	return fmt.Sprintf(fmtString, id.TargetResourceID, id.SettingName)
}

func ParseAdvancedThreatProtectionID(input string) (*AdvancedThreatProtectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Advanced Threat Protection Set ID %q: %+v", input, err)
	}

	parts := strings.Split(input, "/providers/Microsoft.Security/advancedThreatProtectionSettings/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Determining target resource ID, resource ID in unexpected format: %q", id)
	}

	return &AdvancedThreatProtectionId{
		TargetResourceID: parts[0],
		SettingName:      parts[1],
	}, nil
}
