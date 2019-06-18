package datafactory

import (
	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/ar"
)

type Client struct {
	DatasetClient       datafactory.DatasetsClient
	FactoriesClient     datafactory.FactoriesClient
	LinkedServiceClient datafactory.LinkedServicesClient
	PipelinesClient     datafactory.PipelinesClient
}

func BuildClient(endpoint, subscriptionId string, o *ar.ClientOptions) *Client {
	c := Client{}

	c.DatasetClient = datafactory.NewDatasetsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.DatasetClient.Client, o)

	c.FactoriesClient = datafactory.NewFactoriesClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.FactoriesClient.Client, o)

	c.LinkedServiceClient = datafactory.NewLinkedServicesClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.LinkedServiceClient.Client, o)

	c.PipelinesClient = datafactory.NewPipelinesClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.PipelinesClient.Client, o)

	return &c
}
