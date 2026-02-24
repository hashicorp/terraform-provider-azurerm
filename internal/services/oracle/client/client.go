// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	oracle "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	OracleClient *oracle.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	// ORP (backend for ODB@A) partly builds its idempotency keys based on correlationIds sent by client.
	// It seems that AzureRM provider sends the same correlationId for each request during an apply.
	// We need each request to have a different correlationId. By disabling this, Azure will provide a unique correlationId instead.
	tmpClientOptions := *o
	tmpClientOptions.DisableCorrelationRequestID = true
	oracleClient, err := oracle.NewClientWithBaseURI(tmpClientOptions.Environment.ResourceManager, func(c *resourcemanager.Client) {
		tmpClientOptions.Configure(c, tmpClientOptions.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building Database client: %+v", err)
	}

	return &Client{
		OracleClient: oracleClient,
	}, nil
}
