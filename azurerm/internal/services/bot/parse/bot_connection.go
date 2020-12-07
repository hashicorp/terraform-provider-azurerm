package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
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
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Bot Service Name %q", id.BotServiceName),
		fmt.Sprintf("Connection Name %q", id.ConnectionName),
	}
	return strings.Join(segments, " / ")
}

func (id BotConnectionId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.BotService/botServices/%s/connections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.BotServiceName, id.ConnectionName)
}

// BotConnectionID parses a BotConnection ID into an BotConnectionId struct
func BotConnectionID(input string) (*BotConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
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
