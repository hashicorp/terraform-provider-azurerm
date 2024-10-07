// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/voiceservices/2023-04-03/communicationsgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/voiceservices/2023-04-03/testlines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	CommunicationsGatewaysClient *communicationsgateways.CommunicationsGatewaysClient
	TestLinesClient              *testlines.TestLinesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {

	communicationsGatewaysClient, err := communicationsgateways.NewCommunicationsGatewaysClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Communications Gateways Client: %+v", err)
	}
	o.Configure(communicationsGatewaysClient.Client, o.Authorizers.ResourceManager)

	testLinesClient, err := testlines.NewTestLinesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Test Lines Client: %+v", err)
	}
	o.Configure(testLinesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		CommunicationsGatewaysClient: communicationsGatewaysClient,
		TestLinesClient:              testLinesClient,
	}, nil
}
