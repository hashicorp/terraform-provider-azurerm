// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AutomanageConfigurationId struct {
	SubscriptionId           string
	ResourceGroup            string
	ConfigurationProfileName string
}

func NewAutomanageConfigurationID(subscriptionId, resourceGroup, configurationProfileName string) AutomanageConfigurationId {
	return AutomanageConfigurationId{
		SubscriptionId:           subscriptionId,
		ResourceGroup:            resourceGroup,
		ConfigurationProfileName: configurationProfileName,
	}
}

func (id AutomanageConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Configuration Profile Name %q", id.ConfigurationProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Automanage Configuration", segmentsStr)
}

func (id AutomanageConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automanage/configurationProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ConfigurationProfileName)
}

// AutomanageConfigurationID parses a AutomanageConfiguration ID into an AutomanageConfigurationId struct
func AutomanageConfigurationID(input string) (*AutomanageConfigurationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AutomanageConfiguration ID: %+v", input, err)
	}

	resourceId := AutomanageConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ConfigurationProfileName, err = id.PopSegment("configurationProfiles"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
