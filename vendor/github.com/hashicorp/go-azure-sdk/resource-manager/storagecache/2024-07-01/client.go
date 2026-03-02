package v2024_07_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/amlfilesystems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/ascusages"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/autoexportjob"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/autoexportjobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/caches"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/importjobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/skus"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/storagetargets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01/usagemodels"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AmlFilesystems *amlfilesystems.AmlFilesystemsClient
	AscUsages      *ascusages.AscUsagesClient
	AutoExportJob  *autoexportjob.AutoExportJobClient
	AutoExportJobs *autoexportjobs.AutoExportJobsClient
	Caches         *caches.CachesClient
	ImportJobs     *importjobs.ImportJobsClient
	SKUs           *skus.SKUsClient
	StorageTargets *storagetargets.StorageTargetsClient
	UsageModels    *usagemodels.UsageModelsClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	amlFilesystemsClient, err := amlfilesystems.NewAmlFilesystemsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AmlFilesystems client: %+v", err)
	}
	configureFunc(amlFilesystemsClient.Client)

	ascUsagesClient, err := ascusages.NewAscUsagesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AscUsages client: %+v", err)
	}
	configureFunc(ascUsagesClient.Client)

	autoExportJobClient, err := autoexportjob.NewAutoExportJobClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AutoExportJob client: %+v", err)
	}
	configureFunc(autoExportJobClient.Client)

	autoExportJobsClient, err := autoexportjobs.NewAutoExportJobsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AutoExportJobs client: %+v", err)
	}
	configureFunc(autoExportJobsClient.Client)

	cachesClient, err := caches.NewCachesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Caches client: %+v", err)
	}
	configureFunc(cachesClient.Client)

	importJobsClient, err := importjobs.NewImportJobsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ImportJobs client: %+v", err)
	}
	configureFunc(importJobsClient.Client)

	sKUsClient, err := skus.NewSKUsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SKUs client: %+v", err)
	}
	configureFunc(sKUsClient.Client)

	storageTargetsClient, err := storagetargets.NewStorageTargetsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building StorageTargets client: %+v", err)
	}
	configureFunc(storageTargetsClient.Client)

	usageModelsClient, err := usagemodels.NewUsageModelsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building UsageModels client: %+v", err)
	}
	configureFunc(usageModelsClient.Client)

	return &Client{
		AmlFilesystems: amlFilesystemsClient,
		AscUsages:      ascUsagesClient,
		AutoExportJob:  autoExportJobClient,
		AutoExportJobs: autoExportJobsClient,
		Caches:         cachesClient,
		ImportJobs:     importJobsClient,
		SKUs:           sKUsClient,
		StorageTargets: storageTargetsClient,
		UsageModels:    usageModelsClient,
	}, nil
}
