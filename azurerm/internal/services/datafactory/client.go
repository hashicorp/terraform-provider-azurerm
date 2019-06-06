package datafactory

import "github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"

type Client struct {
	DatasetClient       datafactory.DatasetsClient
	FactoriesClient     datafactory.FactoriesClient
	LinkedServiceClient datafactory.LinkedServicesClient
	PipelinesClient     datafactory.PipelinesClient
}
