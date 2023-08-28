// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AutomanageConfigurationVersionId struct {
	SubscriptionId           string
	ResourceGroup            string
	ConfigurationProfileName string
	VersionName              string
}

func NewAutomanageConfigurationVersionID(subscriptionId, resourceGroup, configurationProfileName, versionName string) AutomanageConfigurationVersionId {
	return AutomanageConfigurationVersionId{
		SubscriptionId:           subscriptionId,
		ResourceGroup:            resourceGroup,
		ConfigurationProfileName: configurationProfileName,
		VersionName:              versionName,
	}
}

func (id AutomanageConfigurationVersionId) String() string {
	segments := []string{
		fmt.Sprintf("Version Name %q", id.VersionName),
		fmt.Sprintf("Configuration Profile Name %q", id.ConfigurationProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Automanage Configuration Version", segmentsStr)
}

func (id AutomanageConfigurationVersionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automanage/configurationProfiles/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ConfigurationProfileName, id.VersionName)
}

// AutomanageConfigurationVersionID parses a AutomanageConfigurationVersion ID into an AutomanageConfigurationVersionId struct
func AutomanageConfigurationVersionID(input string) (*AutomanageConfigurationVersionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AutomanageConfigurationVersion ID: %+v", input, err)
	}

	resourceId := AutomanageConfigurationVersionId{
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
	if resourceId.VersionName, err = id.PopSegment("versions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
