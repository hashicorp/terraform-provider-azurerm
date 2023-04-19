package client

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2022-01-01/application"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2022-01-01/batchaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2022-01-01/certificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2022-01-01/pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	batchDataplane "github.com/tombuildsstuff/kermit/sdk/batch/2022-01.15.0/batch"
)

type Client struct {
	AccountClient     *batchaccount.BatchAccountClient
	ApplicationClient *application.ApplicationClient
	CertificateClient *certificate.CertificateClient
	PoolClient        *pool.PoolClient

	BatchManagementAuthorizer autorest.Authorizer
}

func NewClient(o *common.ClientOptions) *Client {
	accountClient := batchaccount.NewBatchAccountClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&accountClient.Client, o.ResourceManagerAuthorizer)

	applicationClient := application.NewApplicationClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&applicationClient.Client, o.ResourceManagerAuthorizer)

	certificateClient := certificate.NewCertificateClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&certificateClient.Client, o.ResourceManagerAuthorizer)

	poolClient := pool.NewPoolClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&poolClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountClient:             &accountClient,
		ApplicationClient:         &applicationClient,
		CertificateClient:         &certificateClient,
		PoolClient:                &poolClient,
		BatchManagementAuthorizer: o.BatchManagementAuthorizer,
	}
}

func (r *Client) JobClient(ctx context.Context, accountId batchaccount.BatchAccountId) (*batchDataplane.JobClient, error) {
	// Retrieve the batch account to find the batch account endpoint
	accountClient := r.AccountClient
	account, err := accountClient.Get(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", accountId, err)
	}
	if model := account.Model; model != nil {
		if model.Properties == nil {
			return nil, fmt.Errorf(`unexpected nil of "AccountProperties" of %s`, accountId)
		}
		if model.Properties.AccountEndpoint == nil {
			return nil, fmt.Errorf(`unexpected nil of "AccountProperties.AccountEndpoint" of %s`, accountId)
		}
	}

	// Copy the client since we'll manipulate its BatchURL
	endpoint := "https://" + *account.Model.Properties.AccountEndpoint
	c := batchDataplane.NewJobClient(endpoint)
	c.BaseClient.Client.Authorizer = r.BatchManagementAuthorizer
	return &c, nil
}
