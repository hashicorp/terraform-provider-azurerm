// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	application "github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/applications"
	batchaccount "github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccounts"
	certificate "github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/certificates"
	pool "github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	batchDataplane "github.com/jackofallops/kermit/sdk/batch/2022-01.15.0/batch"
)

type Client struct {
	AccountClient     *batchaccount.BatchAccountsClient
	ApplicationClient *application.ApplicationsClient
	CertificateClient *certificate.CertificatesClient
	PoolClient        *pool.PoolsClient

	BatchManagementAuthorizer autorest.Authorizer
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	accountClient, err := batchaccount.NewBatchAccountsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Account client: %+v", err)
	}
	o.Configure(accountClient.Client, o.Authorizers.ResourceManager)

	applicationClient, err := application.NewApplicationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Application client: %+v", err)
	}
	o.Configure(applicationClient.Client, o.Authorizers.ResourceManager)

	certificateClient, err := certificate.NewCertificatesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Certificate client: %+v", err)
	}
	o.Configure(certificateClient.Client, o.Authorizers.ResourceManager)

	poolClient, err := pool.NewPoolsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Pool client: %+v", err)
	}
	o.Configure(poolClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AccountClient:             accountClient,
		ApplicationClient:         applicationClient,
		CertificateClient:         certificateClient,
		PoolClient:                poolClient,
		BatchManagementAuthorizer: o.BatchManagementAuthorizer,
	}, nil
}

func (r *Client) JobClient(ctx context.Context, accountId batchaccount.BatchAccountId) (*batchDataplane.JobClient, error) {
	// Retrieve the batch account to find the batch account endpoint
	accountClient := r.AccountClient
	account, err := accountClient.BatchAccountGet(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", accountId, err)
	}

	endpoint := ""
	if account.Model != nil && account.Model.Properties != nil {
		endpoint = "https://" + *account.Model.Properties.AccountEndpoint
	}
	if endpoint == "" {
		return nil, fmt.Errorf("retrieving %s: `properties.AccountEndpoint` was empty", accountId)
	}

	// Copy the client since we'll manipulate its BatchURL
	c := batchDataplane.NewJobClient(endpoint)
	c.Authorizer = r.BatchManagementAuthorizer
	return &c, nil
}
