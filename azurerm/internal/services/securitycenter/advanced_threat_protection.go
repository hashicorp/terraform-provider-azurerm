package securitycenter

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AdvancedThreatProtectionResourceID struct {
	Base azure.ResourceID

	TargetResourceID string
}

func ParseAdvancedThreatProtectionID(input string) (*AdvancedThreatProtectionResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Advanced Threat Protection Set ID %q: %+v", input, err)
	}

	parts := strings.Split(input, "/providers/Microsoft.Security/advancedThreatProtectionSettings/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Error determining target resource ID, resource ID in unexpected format: %q", id)
	}

	return &AdvancedThreatProtectionResourceID{
		Base:             *id,
		TargetResourceID: parts[0],
	}, nil
}
