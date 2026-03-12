// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ServerDNSAliasId struct {
	SubscriptionId string
	ResourceGroup  string
	ServerName     string
	DnsAliaseName  string
}

func NewServerDNSAliasID(subscriptionId, resourceGroup, serverName, dnsAliaseName string) ServerDNSAliasId {
	return ServerDNSAliasId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServerName:     serverName,
		DnsAliaseName:  dnsAliaseName,
	}
}

func (id ServerDNSAliasId) String() string {
	segments := []string{
		fmt.Sprintf("Dns Aliase Name %q", id.DnsAliaseName),
		fmt.Sprintf("Server Name %q", id.ServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "ServerDNS Alias", segmentsStr)
}

func (id ServerDNSAliasId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/dnsAliases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerName, id.DnsAliaseName)
}

// ServerDNSAliasID parses a ServerDNSAlias ID into an ServerDNSAliasId struct
func ServerDNSAliasID(input string) (*ServerDNSAliasId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ServerDNSAlias ID: %+v", input, err)
	}

	resourceId := ServerDNSAliasId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}
	if resourceId.DnsAliaseName, err = id.PopSegment("dnsAliases"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
