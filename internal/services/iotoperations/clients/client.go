// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerauthorization"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

// Client wraps the SDK client(s) used by the iotoperations resources.
type Client struct {
    BrokerAuthorizationClient *brokerauthorization.BrokerAuthorizationClient
}

// NewClient builds the iotoperations clients used by the provider.
func NewClient(o *common.ClientOptions) (*Client, error) {
    brokerAuthorizationClient := brokerauthorization.NewBrokerAuthorizationClientWithBaseURI(o.ResourceManagerEndpoint)
    brokerAuthorizationClient.Client.Authorizer = o.ResourceManagerAuthorizer

    return &Client{
        BrokerAuthorizationClient: brokerAuthorizationClient,
    }, nil
}