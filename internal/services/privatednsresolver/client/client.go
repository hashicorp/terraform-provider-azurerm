// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsforwardingrulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsresolvers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/forwardingrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/inboundendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/outboundendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DnsForwardingRulesetsClient *dnsforwardingrulesets.DnsForwardingRulesetsClient
	DnsResolversClient          *dnsresolvers.DnsResolversClient
	ForwardingRulesClient       *forwardingrules.ForwardingRulesClient
	InboundEndpointsClient      *inboundendpoints.InboundEndpointsClient
	OutboundEndpointsClient     *outboundendpoints.OutboundEndpointsClient
	VirtualNetworkLinksClient   *virtualnetworklinks.VirtualNetworkLinksClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	dnsForwardingRulesetsClient, err := dnsforwardingrulesets.NewDnsForwardingRulesetsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(dnsForwardingRulesetsClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DnsForwardingRulesetsClient client: %+v", err)
	}

	dnsResolversClient, err := dnsresolvers.NewDnsResolversClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(dnsResolversClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DnsResolversClient client: %+v", err)
	}

	forwardingRulesClient, err := forwardingrules.NewForwardingRulesClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(forwardingRulesClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ForwardingRulesClient client: %+v", err)
	}

	inboundEndpointsClient, err := inboundendpoints.NewInboundEndpointsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(inboundEndpointsClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building InboundEndpointsClient client: %+v", err)
	}

	outboundEndpointsClient, err := outboundendpoints.NewOutboundEndpointsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(outboundEndpointsClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building OutboundEndpointsClient client: %+v", err)
	}

	virtualNetworkLinksClient, err := virtualnetworklinks.NewVirtualNetworkLinksClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(virtualNetworkLinksClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VirtualNetworkLinksClient client: %+v", err)
	}

	return &Client{
		DnsForwardingRulesetsClient: dnsForwardingRulesetsClient,
		DnsResolversClient:          dnsResolversClient,
		ForwardingRulesClient:       forwardingRulesClient,
		InboundEndpointsClient:      inboundEndpointsClient,
		OutboundEndpointsClient:     outboundEndpointsClient,
		VirtualNetworkLinksClient:   virtualNetworkLinksClient,
	}, nil
}
