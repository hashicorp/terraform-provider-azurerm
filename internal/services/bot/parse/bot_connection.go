// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type BotConnectionId struct {
	SubscriptionId string
	ResourceGroup  string
	BotServiceName string
	ConnectionName string
}

func NewBotConnectionID(subscriptionId, resourceGroup, botServiceName, connectionName string) BotConnectionId {
	return BotConnectionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		BotServiceName: botServiceName,
		ConnectionName: connectionName,
	}
}

func (id BotConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Connection Name %q", id.ConnectionName),
		fmt.Sprintf("Bot Service Name %q", id.BotServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Bot Connection", segmentsStr)
}

func (id BotConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.BotService/botServices/%s/connections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.BotServiceName, id.ConnectionName)
}

// BotConnectionID parses a BotConnection ID into an BotConnectionId struct
func BotConnectionID(input string) (*BotConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an BotConnection ID: %+v", input, err)
	}

	resourceId := BotConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.BotServiceName, err = id.PopSegment("botServices"); err != nil {
		return nil, err
	}
	if resourceId.ConnectionName, err = id.PopSegment("connections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
