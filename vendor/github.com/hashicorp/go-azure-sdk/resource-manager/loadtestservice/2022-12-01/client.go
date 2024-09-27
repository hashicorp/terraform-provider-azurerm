package v2022_12_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/loadtests"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/quotas"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	LoadTests *loadtests.LoadTestsClient
	Quotas    *quotas.QuotasClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	loadTestsClient, err := loadtests.NewLoadTestsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LoadTests client: %+v", err)
	}
	configureFunc(loadTestsClient.Client)

	quotasClient, err := quotas.NewQuotasClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Quotas client: %+v", err)
	}
	configureFunc(quotasClient.Client)

	return &Client{
		LoadTests: loadTestsClient,
		Quotas:    quotasClient,
	}, nil
}
