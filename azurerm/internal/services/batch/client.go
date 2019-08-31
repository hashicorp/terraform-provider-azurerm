package batch

import (
	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2018-12-01/batch"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AccountClient     *batch.AccountClient
	ApplicationClient *batch.ApplicationClient
	CertificateClient *batch.CertificateClient
	PoolClient        *batch.PoolClient
}

func BuildClient(o *common.ClientOptions) *Client {

	AccountClient := batch.NewAccountClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AccountClient.Client, o.ResourceManagerAuthorizer)

	ApplicationClient := batch.NewApplicationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ApplicationClient.Client, o.ResourceManagerAuthorizer)

	CertificateClient := batch.NewCertificateClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&CertificateClient.Client, o.ResourceManagerAuthorizer)

	PoolClient := batch.NewPoolClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PoolClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountClient:     &AccountClient,
		ApplicationClient: &ApplicationClient,
		CertificateClient: &CertificateClient,
		PoolClient:        &PoolClient,
	}
}
