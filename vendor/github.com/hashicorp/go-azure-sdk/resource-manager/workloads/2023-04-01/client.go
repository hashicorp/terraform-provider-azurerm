package v2023_04_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/providerinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapapplicationserverinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapavailabilityzonedetails"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapcentralinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapdatabaseinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapdiskconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/saplandscapemonitor"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/saprecommendations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapsupportedsku"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapvirtualinstances"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Monitors                      *monitors.MonitorsClient
	ProviderInstances             *providerinstances.ProviderInstancesClient
	SAPApplicationServerInstances *sapapplicationserverinstances.SAPApplicationServerInstancesClient
	SAPAvailabilityZoneDetails    *sapavailabilityzonedetails.SAPAvailabilityZoneDetailsClient
	SAPCentralInstances           *sapcentralinstances.SAPCentralInstancesClient
	SAPDatabaseInstances          *sapdatabaseinstances.SAPDatabaseInstancesClient
	SAPDiskConfigurations         *sapdiskconfigurations.SAPDiskConfigurationsClient
	SAPRecommendations            *saprecommendations.SAPRecommendationsClient
	SAPSupportedSku               *sapsupportedsku.SAPSupportedSkuClient
	SAPVirtualInstances           *sapvirtualinstances.SAPVirtualInstancesClient
	SapLandscapeMonitor           *saplandscapemonitor.SapLandscapeMonitorClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	monitorsClient, err := monitors.NewMonitorsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Monitors client: %+v", err)
	}
	configureFunc(monitorsClient.Client)

	providerInstancesClient, err := providerinstances.NewProviderInstancesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ProviderInstances client: %+v", err)
	}
	configureFunc(providerInstancesClient.Client)

	sAPApplicationServerInstancesClient, err := sapapplicationserverinstances.NewSAPApplicationServerInstancesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SAPApplicationServerInstances client: %+v", err)
	}
	configureFunc(sAPApplicationServerInstancesClient.Client)

	sAPAvailabilityZoneDetailsClient, err := sapavailabilityzonedetails.NewSAPAvailabilityZoneDetailsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SAPAvailabilityZoneDetails client: %+v", err)
	}
	configureFunc(sAPAvailabilityZoneDetailsClient.Client)

	sAPCentralInstancesClient, err := sapcentralinstances.NewSAPCentralInstancesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SAPCentralInstances client: %+v", err)
	}
	configureFunc(sAPCentralInstancesClient.Client)

	sAPDatabaseInstancesClient, err := sapdatabaseinstances.NewSAPDatabaseInstancesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SAPDatabaseInstances client: %+v", err)
	}
	configureFunc(sAPDatabaseInstancesClient.Client)

	sAPDiskConfigurationsClient, err := sapdiskconfigurations.NewSAPDiskConfigurationsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SAPDiskConfigurations client: %+v", err)
	}
	configureFunc(sAPDiskConfigurationsClient.Client)

	sAPRecommendationsClient, err := saprecommendations.NewSAPRecommendationsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SAPRecommendations client: %+v", err)
	}
	configureFunc(sAPRecommendationsClient.Client)

	sAPSupportedSkuClient, err := sapsupportedsku.NewSAPSupportedSkuClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SAPSupportedSku client: %+v", err)
	}
	configureFunc(sAPSupportedSkuClient.Client)

	sAPVirtualInstancesClient, err := sapvirtualinstances.NewSAPVirtualInstancesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SAPVirtualInstances client: %+v", err)
	}
	configureFunc(sAPVirtualInstancesClient.Client)

	sapLandscapeMonitorClient, err := saplandscapemonitor.NewSapLandscapeMonitorClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SapLandscapeMonitor client: %+v", err)
	}
	configureFunc(sapLandscapeMonitorClient.Client)

	return &Client{
		Monitors:                      monitorsClient,
		ProviderInstances:             providerInstancesClient,
		SAPApplicationServerInstances: sAPApplicationServerInstancesClient,
		SAPAvailabilityZoneDetails:    sAPAvailabilityZoneDetailsClient,
		SAPCentralInstances:           sAPCentralInstancesClient,
		SAPDatabaseInstances:          sAPDatabaseInstancesClient,
		SAPDiskConfigurations:         sAPDiskConfigurationsClient,
		SAPRecommendations:            sAPRecommendationsClient,
		SAPSupportedSku:               sAPSupportedSkuClient,
		SAPVirtualInstances:           sAPVirtualInstancesClient,
		SapLandscapeMonitor:           sapLandscapeMonitorClient,
	}, nil
}
