package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2020-01-01/getrecommendations" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	RecommendationsClient *getrecommendations.GetRecommendationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	recommendationsClient := getrecommendations.NewGetRecommendationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&recommendationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		RecommendationsClient: &recommendationsClient,
	}
}
