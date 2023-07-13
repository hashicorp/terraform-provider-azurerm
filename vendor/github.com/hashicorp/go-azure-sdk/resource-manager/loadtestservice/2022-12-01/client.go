package v2022_12_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/loadtest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/loadtests"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/quotas"
)

type Client struct {
	LoadTest  *loadtest.LoadTestClient
	LoadTests *loadtests.LoadTestsClient
	Quotas    *quotas.QuotasClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	loadTestClient := loadtest.NewLoadTestClientWithBaseURI(endpoint)
	configureAuthFunc(&loadTestClient.Client)

	loadTestsClient := loadtests.NewLoadTestsClientWithBaseURI(endpoint)
	configureAuthFunc(&loadTestsClient.Client)

	quotasClient := quotas.NewQuotasClientWithBaseURI(endpoint)
	configureAuthFunc(&quotasClient.Client)

	return Client{
		LoadTest:  &loadTestClient,
		LoadTests: &loadTestsClient,
		Quotas:    &quotasClient,
	}
}
