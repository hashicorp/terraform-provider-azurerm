// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

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

func NewClient(o *common.ClientOptions) (*Client, error) {
	accountClient, err := account.NewAccountClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Account client: %+v", err)
	}
	o.Configure(accountClient.Client, o.Authorizers.ResourceManager)

	dataSetClient, err := dataset.NewDataSetClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DataSet client: %+v", err)
	}
	o.Configure(dataSetClient.Client, o.Authorizers.ResourceManager)

	sharesClient, err := share.NewShareClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Share client: %+v", err)
	}
	o.Configure(sharesClient.Client, o.Authorizers.ResourceManager)

	synchronizationSettingsClient, err := synchronizationsetting.NewSynchronizationSettingClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SynchronizationSetting client: %+v", err)
	}
	o.Configure(synchronizationSettingsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AccountClient:         accountClient,
		DataSetClient:         dataSetClient,
		SharesClient:          sharesClient,
		SynchronizationClient: synchronizationSettingsClient,
	}, nil
}
