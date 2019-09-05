package datafactory

import (
	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DatasetClient       *datafactory.DatasetsClient
	FactoriesClient     *datafactory.FactoriesClient
	LinkedServiceClient *datafactory.LinkedServicesClient
	PipelinesClient     *datafactory.PipelinesClient
}

func BuildClient(o *common.ClientOptions) *Client {

	DatasetClient := datafactory.NewDatasetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatasetClient.Client, o.ResourceManagerAuthorizer)

	FactoriesClient := datafactory.NewFactoriesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&FactoriesClient.Client, o.ResourceManagerAuthorizer)

	LinkedServiceClient := datafactory.NewLinkedServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&LinkedServiceClient.Client, o.ResourceManagerAuthorizer)

	PipelinesClient := datafactory.NewPipelinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PipelinesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DatasetClient:       &DatasetClient,
		FactoriesClient:     &FactoriesClient,
		LinkedServiceClient: &LinkedServiceClient,
		PipelinesClient:     &PipelinesClient,
	}
}
