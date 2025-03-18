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

type FirewallNetworkRuleCollectionId struct {
	SubscriptionId            string
	ResourceGroup             string
	AzureFirewallName         string
	NetworkRuleCollectionName string
}

func NewFirewallNetworkRuleCollectionID(subscriptionId, resourceGroup, azureFirewallName, networkRuleCollectionName string) FirewallNetworkRuleCollectionId {
	return FirewallNetworkRuleCollectionId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		AzureFirewallName:         azureFirewallName,
		NetworkRuleCollectionName: networkRuleCollectionName,
	}
}

func (id FirewallNetworkRuleCollectionId) String() string {
	segments := []string{
		fmt.Sprintf("Network Rule Collection Name %q", id.NetworkRuleCollectionName),
		fmt.Sprintf("Azure Firewall Name %q", id.AzureFirewallName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Firewall Network Rule Collection", segmentsStr)
}

func (id FirewallNetworkRuleCollectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/azureFirewalls/%s/networkRuleCollections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AzureFirewallName, id.NetworkRuleCollectionName)
}

// FirewallNetworkRuleCollectionID parses a FirewallNetworkRuleCollection ID into an FirewallNetworkRuleCollectionId struct
func FirewallNetworkRuleCollectionID(input string) (*FirewallNetworkRuleCollectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an FirewallNetworkRuleCollection ID: %+v", input, err)
	}

	resourceId := FirewallNetworkRuleCollectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.AzureFirewallName, err = id.PopSegment("azureFirewalls"); err != nil {
		return nil, err
	}
	if resourceId.NetworkRuleCollectionName, err = id.PopSegment("networkRuleCollections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
