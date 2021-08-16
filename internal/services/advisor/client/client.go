package client

import (
	"github.com/Azure/azure-sdk-for-go/services/advisor/mgmt/2020-01-01/advisor"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	RecommendationsClient *advisor.RecommendationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	recommendationsClient := advisor.NewRecommendationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&recommendationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		RecommendationsClient: &recommendationsClient,
	}
}
