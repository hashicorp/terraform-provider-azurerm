package datafactory

import (
	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DatasetClient       datafactory.DatasetsClient
	FactoriesClient     datafactory.FactoriesClient
	LinkedServiceClient datafactory.LinkedServicesClient
	PipelinesClient     datafactory.PipelinesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.DatasetClient = datafactory.NewDatasetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DatasetClient.Client, o.ResourceManagerAuthorizer)

	c.FactoriesClient = datafactory.NewFactoriesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.FactoriesClient.Client, o.ResourceManagerAuthorizer)

	c.LinkedServiceClient = datafactory.NewLinkedServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.LinkedServiceClient.Client, o.ResourceManagerAuthorizer)

	c.PipelinesClient = datafactory.NewPipelinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.PipelinesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
