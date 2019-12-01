package clients

import (
	"github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2017-08-01/analysisservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type AnalysisServicesClient struct {
	ServerClient *analysisservices.ServersClient
}

func newAnalysisServicesClient(o *common.ClientOptions) *AnalysisServicesClient {
	serverClient := analysisservices.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverClient.Client, o.ResourceManagerAuthorizer)

	return &AnalysisServicesClient{
		ServerClient: &serverClient,
	}
}
