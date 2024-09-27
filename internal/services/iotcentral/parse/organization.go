// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type OrganizationId struct {
	SubscriptionId string
	ResourceGroup  string
	IotAppName     string
	Name           string
}

func NewOrganizationID(subscriptionId, resourceGroup, iotAppName, name string) OrganizationId {
	return OrganizationId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		IotAppName:     iotAppName,
		Name:           name,
	}
}

func (id OrganizationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Iot App Name %q", id.IotAppName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Organization", segmentsStr)
}

func (id OrganizationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.IoTCentral/iotApps/%s/organizations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IotAppName, id.Name)
}

// OrganizationID parses a Organization ID into an OrganizationId struct
func OrganizationID(input string) (*OrganizationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an Organization ID: %+v", input, err)
	}

	resourceId := OrganizationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.IotAppName, err = id.PopSegment("iotApps"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("organizations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
