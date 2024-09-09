// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package client

import (
	oracedatabase "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	OracleDatabaseClient *oracedatabase.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	// ORP (backend for ODB@A) partly builds its idempotency keys based on correlationIds sent by client.
	// It seems that AzureRM provider sends the same correlationId for each request during an apply.
	// We need each request to have a different correlationId. By disabling this, Azure will provide a unique correlationId instead.
	o.DisableCorrelationRequestID = true
	oracleDatabaseClient, err := oracedatabase.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		OracleDatabaseClient: oracleDatabaseClient,
	}, nil
}
