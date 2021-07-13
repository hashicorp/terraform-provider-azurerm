package client

import (
	"github.com/Azure/azure-sdk-for-go/services/batch/2020-03-01.11.0/batchDataplane"
	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2020-03-01/batch"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AccountClient     *batch.AccountClient
	ApplicationClient *batch.ApplicationClient
	CertificateClient *batch.CertificateClient
	PoolClient        *batch.PoolClient

	JobClient *batchDataplane.JobClient
}

func NewClient(o *common.ClientOptions) *Client {
	accountClient := batch.NewAccountClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&accountClient.Client, o.ResourceManagerAuthorizer)

	applicationClient := batch.NewApplicationClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&applicationClient.Client, o.ResourceManagerAuthorizer)

	certificateClient := batch.NewCertificateClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&certificateClient.Client, o.ResourceManagerAuthorizer)

	poolClient := batch.NewPoolClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&poolClient.Client, o.ResourceManagerAuthorizer)

	jobClient := batchDataplane.NewJobClient("")
	o.ConfigureClient(&jobClient.Client, o.BatchAuthorizer)

	return &Client{
		AccountClient:     &accountClient,
		ApplicationClient: &applicationClient,
		CertificateClient: &certificateClient,
		PoolClient:        &poolClient,
		JobClient:         &jobClient,
	}
}
