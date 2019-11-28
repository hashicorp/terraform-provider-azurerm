package domainservices

import (
	"github.com/Azure/azure-sdk-for-go/services/domainservices/mgmt/2017-06-01/aad"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DomainServicesClient *aad.DomainServicesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	DomainServicesClient := aad.NewDomainServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DomainServicesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DomainServicesClient: &DomainServicesClient,
	}
}
