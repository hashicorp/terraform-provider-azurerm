package batch

import (
	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2018-12-01/batch"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AccountClient     batch.AccountClient
	ApplicationClient batch.ApplicationClient
	CertificateClient batch.CertificateClient
	PoolClient        batch.PoolClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.AccountClient = batch.NewAccountClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AccountClient.Client, o.ResourceManagerAuthorizer)

	c.ApplicationClient = batch.NewApplicationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ApplicationClient.Client, o.ResourceManagerAuthorizer)

	c.CertificateClient = batch.NewCertificateClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.CertificateClient.Client, o.ResourceManagerAuthorizer)

	c.PoolClient = batch.NewPoolClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.PoolClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
