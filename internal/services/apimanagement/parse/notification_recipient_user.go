// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type NotificationRecipientUserId struct {
	SubscriptionId    string
	ResourceGroup     string
	ServiceName       string
	NotificationName  string
	RecipientUserName string
}

func NewNotificationRecipientUserID(subscriptionId, resourceGroup, serviceName, notificationName, recipientUserName string) NotificationRecipientUserId {
	return NotificationRecipientUserId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		ServiceName:       serviceName,
		NotificationName:  notificationName,
		RecipientUserName: recipientUserName,
	}
}

func (id NotificationRecipientUserId) String() string {
	segments := []string{
		fmt.Sprintf("Recipient User Name %q", id.RecipientUserName),
		fmt.Sprintf("Notification Name %q", id.NotificationName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Notification Recipient User", segmentsStr)
}

func (id NotificationRecipientUserId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/notifications/%s/recipientUsers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.NotificationName, id.RecipientUserName)
}

// NotificationRecipientUserID parses a NotificationRecipientUser ID into an NotificationRecipientUserId struct
func NotificationRecipientUserID(input string) (*NotificationRecipientUserId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an NotificationRecipientUser ID: %+v", input, err)
	}

	resourceId := NotificationRecipientUserId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}
	if resourceId.NotificationName, err = id.PopSegment("notifications"); err != nil {
		return nil, err
	}
	if resourceId.RecipientUserName, err = id.PopSegment("recipientUsers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
