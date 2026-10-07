// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/data-plane/batch/2022-01-01-15-0/jobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/application"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/certificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	// Resource Manager
	AccountClient     *batchaccount.BatchAccountClient
	ApplicationClient *application.ApplicationClient
	CertificateClient *certificate.CertificateClient
	PoolClient        *pool.PoolClient

	// Data Plane
	JobsClient *jobs.JobsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	// Resource Manager
	accountClient, err := batchaccount.NewBatchAccountClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Account client: %+v", err)
	}
	o.Configure(accountClient.Client, o.Authorizers.ResourceManager)

	applicationClient, err := application.NewApplicationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Application client: %+v", err)
	}
	o.Configure(applicationClient.Client, o.Authorizers.ResourceManager)

	certificateClient, err := certificate.NewCertificateClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Certificate client: %+v", err)
	}
	o.Configure(certificateClient.Client, o.Authorizers.ResourceManager)

	poolClient, err := pool.NewPoolClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Pool client: %+v", err)
	}
	o.Configure(poolClient.Client, o.Authorizers.ResourceManager)

	// Data Plane
	jobsClient, err := jobs.NewJobsClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Jobs client: %+v", err)
	}
	o.Configure(jobsClient.Client, o.Authorizers.BatchManagement)

	return &Client{
		AccountClient:     accountClient,
		ApplicationClient: applicationClient,
		CertificateClient: certificateClient,
		PoolClient:        poolClient,

		JobsClient: jobsClient,
	}, nil
}
