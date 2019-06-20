package cognitive

import (
	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2017-04-18/cognitiveservices"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AccountsClient cognitiveservices.AccountsClient
}

func BuildClient(endpoint string, authorizer autorest.Authorizer, o *common.ClientOptions) *Client {
	c := Client{}

	c.AccountsClient = cognitiveservices.NewAccountsClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AccountsClient.Client, authorizer)

	return &c
}
