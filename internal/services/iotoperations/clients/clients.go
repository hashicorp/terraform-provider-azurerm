// Copyright (c) HashiCorp,
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/broker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

// Client wraps the SDK client(s) used by the iotoperations resources.
type Client struct {
    BrokerClient *broker.BrokerClient
}

// NewClient builds the iotoperations clients used by the provider.
func NewClient(o *common.ClientOptions) (*Client, error) {
    brokerClient := broker.NewBrokerClientWithBaseURI(o.ResourceManagerEndpoint)
    brokerClient.Client.Authorizer = o.ResourceManagerAuthorizer

    return &Client{
        BrokerClient: brokerClient,
    }, nil
}