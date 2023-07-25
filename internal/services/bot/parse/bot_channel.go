// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type BotChannelId struct {
	SubscriptionId string
	ResourceGroup  string
	BotServiceName string
	ChannelName    string
}

func NewBotChannelID(subscriptionId, resourceGroup, botServiceName, channelName string) BotChannelId {
	return BotChannelId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		BotServiceName: botServiceName,
		ChannelName:    channelName,
	}
}

func (id BotChannelId) String() string {
	segments := []string{
		fmt.Sprintf("Channel Name %q", id.ChannelName),
		fmt.Sprintf("Bot Service Name %q", id.BotServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Bot Channel", segmentsStr)
}

func (id BotChannelId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.BotService/botServices/%s/channels/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.BotServiceName, id.ChannelName)
}

// BotChannelID parses a BotChannel ID into an BotChannelId struct
func BotChannelID(input string) (*BotChannelId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an BotChannel ID: %+v", input, err)
	}

	resourceId := BotChannelId{
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
	if resourceId.ChannelName, err = id.PopSegment("channels"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
