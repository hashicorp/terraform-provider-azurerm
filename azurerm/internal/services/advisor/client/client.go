package client

import (
	"github.com/Azure/azure-sdk-for-go/services/advisor/mgmt/2020-01-01/advisor"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	RecommendationsClient *advisor.RecommendationsClient
	SuppressionsClient    *advisor.SuppressionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	recommendationsClient := advisor.NewRecommendationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&recommendationsClient.Client, o.ResourceManagerAuthorizer)

	suppressionsClient := advisor.NewSuppressionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&suppressionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		RecommendationsClient: &recommendationsClient,
		SuppressionsClient:    &suppressionsClient,
	}
}
