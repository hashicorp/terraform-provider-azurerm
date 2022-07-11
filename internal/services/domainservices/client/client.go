package client

import (
	"github.com/Azure/azure-sdk-for-go/services/domainservices/mgmt/2020-01-01/aad"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DomainServicesClient *aad.DomainServicesClient
}

func NewClient(o *common.ClientOptions) *Client {
	domainServicesClient := aad.NewDomainServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&domainServicesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DomainServicesClient: &domainServicesClient,
	}
}
