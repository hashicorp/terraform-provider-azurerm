// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type NotificationRecipientEmailId struct {
	SubscriptionId     string
	ResourceGroup      string
	ServiceName        string
	NotificationName   string
	RecipientEmailName string
}

func NewNotificationRecipientEmailID(subscriptionId, resourceGroup, serviceName, notificationName, recipientEmailName string) NotificationRecipientEmailId {
	return NotificationRecipientEmailId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		ServiceName:        serviceName,
		NotificationName:   notificationName,
		RecipientEmailName: recipientEmailName,
	}
}

func (id NotificationRecipientEmailId) String() string {
	segments := []string{
		fmt.Sprintf("Recipient Email Name %q", id.RecipientEmailName),
		fmt.Sprintf("Notification Name %q", id.NotificationName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Notification Recipient Email", segmentsStr)
}

func (id NotificationRecipientEmailId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/notifications/%s/recipientEmails/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.NotificationName, id.RecipientEmailName)
}

// NotificationRecipientEmailID parses a NotificationRecipientEmail ID into an NotificationRecipientEmailId struct
func NotificationRecipientEmailID(input string) (*NotificationRecipientEmailId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an NotificationRecipientEmail ID: %+v", input, err)
	}

	resourceId := NotificationRecipientEmailId{
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
	if resourceId.RecipientEmailName, err = id.PopSegment("recipientEmails"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
