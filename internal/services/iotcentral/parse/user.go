// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type UserId struct {
	SubscriptionId string
	ResourceGroup  string
	IotAppName     string
	Name           string
}

func NewUserID(subscriptionId, resourceGroup, iotAppName, name string) UserId {
	return UserId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		IotAppName:     iotAppName,
		Name:           name,
	}
}

func (id UserId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Iot App Name %q", id.IotAppName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "User", segmentsStr)
}

func (id UserId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.IoTCentral/iotApps/%s/users/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IotAppName, id.Name)
}

// UserID parses a User ID into an UserId struct
func UserID(input string) (*UserId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an User ID: %+v", input, err)
	}

	resourceId := UserId{
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
	if resourceId.Name, err = id.PopSegment("users"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
