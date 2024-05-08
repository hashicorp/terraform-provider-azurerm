package v2024_05_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2024-05-01/hdinsights"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Hdinsights *hdinsights.HdinsightsClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	hdinsightsClient, err := hdinsights.NewHdinsightsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Hdinsights client: %+v", err)
	}
	configureFunc(hdinsightsClient.Client)

	return &Client{
		Hdinsights: hdinsightsClient,
	}, nil
}
