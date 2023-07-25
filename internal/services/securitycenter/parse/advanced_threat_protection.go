// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

var _ resourceids.Id = AdvancedThreatProtectionId{}

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

func (id AdvancedThreatProtectionId) ID() string {
	fmtString := "%s/providers/Microsoft.Security/advancedThreatProtectionSettings/%s"
	return fmt.Sprintf(fmtString, id.TargetResourceID, id.SettingName)
}

func (id AdvancedThreatProtectionId) String() string {
	components := []string{
		fmt.Sprintf("Target Resource ID %q", id.TargetResourceID),
		fmt.Sprintf("Setting Name %q", id.SettingName),
	}
	return fmt.Sprintf("Advanced Protection %s", strings.Join(components, " / "))
}

func AdvancedThreatProtectionID(input string) (*AdvancedThreatProtectionId, error) {
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
