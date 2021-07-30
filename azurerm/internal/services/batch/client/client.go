package client

import (
	"context"
	"fmt"

	batchDataplane "github.com/Azure/azure-sdk-for-go/services/batch/2020-03-01.11.0/batch"
	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2020-03-01/batch"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/parse"
)

type Client struct {
	AccountClient     *batch.AccountClient
	ApplicationClient *batch.ApplicationClient
	CertificateClient *batch.CertificateClient
	PoolClient        *batch.PoolClient

	BatchManagementAuthorizer autorest.Authorizer
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

	return &Client{
		AccountClient:             &accountClient,
		ApplicationClient:         &applicationClient,
		CertificateClient:         &certificateClient,
		PoolClient:                &poolClient,
		BatchManagementAuthorizer: o.BatchManagementAuthorizer,
	}
}

func (r *Client) JobClient(ctx context.Context, accountId parse.AccountId) (*batchDataplane.JobClient, error) {
	// Retrieve the batch account to find the batch account endpoint
	accountClient := r.AccountClient
	account, err := accountClient.Get(ctx, accountId.ResourceGroup, accountId.BatchAccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", accountId, err)
	}
	if account.AccountProperties == nil {
		return nil, fmt.Errorf(`unexpected nil of "AccountProperties" of %s`, accountId)
	}
	if account.AccountProperties.AccountEndpoint == nil {
		return nil, fmt.Errorf(`unexpected nil of "AccountProperties.AccountEndpoint" of %s`, accountId)
	}

	// Copy the client since we'll manipulate its BatchURL
	endpoint := "https://" + *account.AccountProperties.AccountEndpoint
	c := batchDataplane.NewJobClient(endpoint)
	c.BaseClient.Client.Authorizer = r.BatchManagementAuthorizer
	return &c, nil
}
