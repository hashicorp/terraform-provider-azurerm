package clients

import (
	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2017-04-18/cognitiveservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type CognitiveServicesClient struct {
	AccountsClient *cognitiveservices.AccountsClient
}

func newCognitiveServicesClient(o *common.ClientOptions) *CognitiveServicesClient {
	accountsClient := cognitiveservices.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&accountsClient.Client, o.ResourceManagerAuthorizer)

	return &CognitiveServicesClient{
		AccountsClient: &accountsClient,
	}
}
