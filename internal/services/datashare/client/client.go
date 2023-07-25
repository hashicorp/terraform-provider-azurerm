// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/account"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/synchronizationsetting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccountClient         *account.AccountClient
	DataSetClient         *dataset.DataSetClient
	SharesClient          *share.ShareClient
	SynchronizationClient *synchronizationsetting.SynchronizationSettingClient
}

func NewClient(o *common.ClientOptions) *Client {
	accountClient := account.NewAccountClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&accountClient.Client, o.ResourceManagerAuthorizer)

	dataSetClient := dataset.NewDataSetClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&dataSetClient.Client, o.ResourceManagerAuthorizer)

	sharesClient := share.NewShareClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&sharesClient.Client, o.ResourceManagerAuthorizer)

	synchronizationSettingsClient := synchronizationsetting.NewSynchronizationSettingClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&synchronizationSettingsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountClient:         &accountClient,
		DataSetClient:         &dataSetClient,
		SharesClient:          &sharesClient,
		SynchronizationClient: &synchronizationSettingsClient,
	}
}
