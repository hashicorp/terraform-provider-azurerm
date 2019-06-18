package cognitive

import (
	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2017-04-18/cognitiveservices"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/ar"
)

type Client struct {
	AccountsClient cognitiveservices.AccountsClient
}

func BuildClient(endpoint, subscriptionId, partnerId string, auth autorest.Authorizer, skipProviderReg bool) *Client {
	c := Client{}

	c.AccountsClient = cognitiveservices.NewAccountsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.AccountsClient.Client, auth, partnerId, skipProviderReg)

	return &c
}
