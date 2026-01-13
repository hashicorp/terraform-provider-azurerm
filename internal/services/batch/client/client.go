// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/applications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/certificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	batchDataplane "github.com/jackofallops/kermit/sdk/batch/2022-01.15.0/batch"
)

type Client struct {
	AccountClient     *batchaccounts.BatchAccountsClient
	ApplicationClient *applications.ApplicationsClient
	CertificateClient *certificates.CertificatesClient
	PoolClient        *pools.PoolsClient

	BatchManagementAuthorizer autorest.Authorizer
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	accountClient, err := batchaccounts.NewBatchAccountsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Account client: %+v", err)
	}
	o.Configure(accountClient.Client, o.Authorizers.ResourceManager)

	applicationClient, err := applications.NewApplicationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Application client: %+v", err)
	}
	o.Configure(applicationClient.Client, o.Authorizers.ResourceManager)

	certificateClient, err := certificates.NewCertificatesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Certificate client: %+v", err)
	}
	o.Configure(certificateClient.Client, o.Authorizers.ResourceManager)

	poolClient, err := pools.NewPoolsClientWithBaseURI(o.Environment.ResourceManager)
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

func (r *Client) JobClient(ctx context.Context, accountId batchaccounts.BatchAccountId) (*batchDataplane.JobClient, error) {
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
