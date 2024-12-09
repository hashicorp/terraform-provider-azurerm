// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type BotHealthbotId struct {
	SubscriptionId string
	ResourceGroup  string
	HealthBotName  string
}

func NewBotHealthbotID(subscriptionId, resourceGroup, healthBotName string) BotHealthbotId {
	return BotHealthbotId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		HealthBotName:  healthBotName,
	}
}

func (id BotHealthbotId) String() string {
	segments := []string{
		fmt.Sprintf("Health Bot Name %q", id.HealthBotName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Bot Healthbot", segmentsStr)
}

func (id BotHealthbotId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HealthBot/healthBots/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.HealthBotName)
}

// BotHealthbotID parses a BotHealthbot ID into an BotHealthbotId struct
func BotHealthbotID(input string) (*BotHealthbotId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an BotHealthbot ID: %+v", input, err)
	}

	resourceId := BotHealthbotId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.HealthBotName, err = id.PopSegment("healthBots"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
