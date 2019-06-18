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

	return &c
}
