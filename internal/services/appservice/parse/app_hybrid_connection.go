// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AppHybridConnectionId struct {
	SubscriptionId                string
	ResourceGroup                 string
	SiteName                      string
	HybridConnectionNamespaceName string
	RelayName                     string
}

func NewAppHybridConnectionID(subscriptionId, resourceGroup, siteName, hybridConnectionNamespaceName, relayName string) AppHybridConnectionId {
	return AppHybridConnectionId{
		SubscriptionId:                subscriptionId,
		ResourceGroup:                 resourceGroup,
		SiteName:                      siteName,
		HybridConnectionNamespaceName: hybridConnectionNamespaceName,
		RelayName:                     relayName,
	}
}

func (id AppHybridConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Relay Name %q", id.RelayName),
		fmt.Sprintf("Hybrid Connection Namespace Name %q", id.HybridConnectionNamespaceName),
		fmt.Sprintf("Site Name %q", id.SiteName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "App Hybrid Connection", segmentsStr)
}

func (id AppHybridConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/hybridConnectionNamespaces/%s/relays/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName)
}

// AppHybridConnectionID parses a AppHybridConnection ID into an AppHybridConnectionId struct
func AppHybridConnectionID(input string) (*AppHybridConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AppHybridConnection ID: %+v", input, err)
	}

	resourceId := AppHybridConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SiteName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}
	if resourceId.HybridConnectionNamespaceName, err = id.PopSegment("hybridConnectionNamespaces"); err != nil {
		return nil, err
	}
	if resourceId.RelayName, err = id.PopSegment("relays"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
