package v2025_09_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagediscovery/2025-09-01/report"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagediscovery/2025-09-01/storagediscoveryworkspaces"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Report                     *report.ReportClient
	StorageDiscoveryWorkspaces *storagediscoveryworkspaces.StorageDiscoveryWorkspacesClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	reportClient, err := report.NewReportClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Report client: %+v", err)
	}
	configureFunc(reportClient.Client)

	storageDiscoveryWorkspacesClient, err := storagediscoveryworkspaces.NewStorageDiscoveryWorkspacesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building StorageDiscoveryWorkspaces client: %+v", err)
	}
	configureFunc(storageDiscoveryWorkspacesClient.Client)

	return &Client{
		Report:                     reportClient,
		StorageDiscoveryWorkspaces: storageDiscoveryWorkspacesClient,
	}, nil
}
