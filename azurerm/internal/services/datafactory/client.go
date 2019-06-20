package datafactory

import (
	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DatasetClient       datafactory.DatasetsClient
	FactoriesClient     datafactory.FactoriesClient
	LinkedServiceClient datafactory.LinkedServicesClient
	PipelinesClient     datafactory.PipelinesClient
}

func BuildClient(endpoint string, authorizer autorest.Authorizer, o *common.ClientOptions) *Client {
	c := Client{}

	c.DatasetClient = datafactory.NewDatasetsClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DatasetClient.Client, authorizer)

	c.FactoriesClient = datafactory.NewFactoriesClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.FactoriesClient.Client, authorizer)

	c.LinkedServiceClient = datafactory.NewLinkedServicesClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.LinkedServiceClient.Client, authorizer)

	c.PipelinesClient = datafactory.NewPipelinesClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.PipelinesClient.Client, authorizer)

	return &c
}
