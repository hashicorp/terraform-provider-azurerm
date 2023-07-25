// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AutoProvisioningSettingId struct {
	SubscriptionId string
	Name           string
}

func NewAutoProvisioningSettingID(subscriptionId, name string) AutoProvisioningSettingId {
	return AutoProvisioningSettingId{
		SubscriptionId: subscriptionId,
		Name:           name,
	}
}

func (id AutoProvisioningSettingId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Auto Provisioning Setting", segmentsStr)
}

func (id AutoProvisioningSettingId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Security/autoProvisioningSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.Name)
}

// AutoProvisioningSettingID parses a AutoProvisioningSetting ID into an AutoProvisioningSettingId struct
func AutoProvisioningSettingID(input string) (*AutoProvisioningSettingId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AutoProvisioningSetting ID: %+v", input, err)
	}

	resourceId := AutoProvisioningSettingId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.Name, err = id.PopSegment("autoProvisioningSettings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// AutoProvisioningSettingIDInsensitively parses an AutoProvisioningSetting ID into an AutoProvisioningSettingId struct, insensitively
// This should only be used to parse an ID for rewriting, the AutoProvisioningSettingID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func AutoProvisioningSettingIDInsensitively(input string) (*AutoProvisioningSettingId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AutoProvisioningSettingId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	// find the correct casing for the 'autoProvisioningSettings' segment
	autoProvisioningSettingsKey := "autoProvisioningSettings"
	for key := range id.Path {
		if strings.EqualFold(key, autoProvisioningSettingsKey) {
			autoProvisioningSettingsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(autoProvisioningSettingsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
