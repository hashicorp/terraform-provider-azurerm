package client

import (
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

func NewClient(o *common.ClientOptions) *Client {

	dnsForwardingRulesetsClient := dnsforwardingrulesets.NewDnsForwardingRulesetsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&dnsForwardingRulesetsClient.Client, o.ResourceManagerAuthorizer)

	dnsResolversClient := dnsresolvers.NewDnsResolversClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&dnsResolversClient.Client, o.ResourceManagerAuthorizer)

	forwardingRulesClient := forwardingrules.NewForwardingRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&forwardingRulesClient.Client, o.ResourceManagerAuthorizer)

	inboundEndpointsClient := inboundendpoints.NewInboundEndpointsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&inboundEndpointsClient.Client, o.ResourceManagerAuthorizer)

	outboundEndpointsClient := outboundendpoints.NewOutboundEndpointsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&outboundEndpointsClient.Client, o.ResourceManagerAuthorizer)

	virtualNetworkLinksClient := virtualnetworklinks.NewVirtualNetworkLinksClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&virtualNetworkLinksClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DnsForwardingRulesetsClient: &dnsForwardingRulesetsClient,
		DnsResolversClient:          &dnsResolversClient,
		ForwardingRulesClient:       &forwardingRulesClient,
		InboundEndpointsClient:      &inboundEndpointsClient,
		OutboundEndpointsClient:     &outboundEndpointsClient,
		VirtualNetworkLinksClient:   &virtualNetworkLinksClient,
	}
}
