package v2023_01_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/elasticsan"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/elasticsans"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/elasticsanskus"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/privatelinkresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumes"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	ElasticSan                 *elasticsan.ElasticSanClient
	ElasticSanSkus             *elasticsanskus.ElasticSanSkusClient
	ElasticSans                *elasticsans.ElasticSansClient
	PrivateEndpointConnections *privateendpointconnections.PrivateEndpointConnectionsClient
	PrivateLinkResources       *privatelinkresources.PrivateLinkResourcesClient
	Snapshots                  *snapshots.SnapshotsClient
	VolumeGroups               *volumegroups.VolumeGroupsClient
	Volumes                    *volumes.VolumesClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	elasticSanClient, err := elasticsan.NewElasticSanClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ElasticSan client: %+v", err)
	}
	configureFunc(elasticSanClient.Client)

	elasticSanSkusClient, err := elasticsanskus.NewElasticSanSkusClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ElasticSanSkus client: %+v", err)
	}
	configureFunc(elasticSanSkusClient.Client)

	elasticSansClient, err := elasticsans.NewElasticSansClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ElasticSans client: %+v", err)
	}
	configureFunc(elasticSansClient.Client)

	privateEndpointConnectionsClient, err := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrivateEndpointConnections client: %+v", err)
	}
	configureFunc(privateEndpointConnectionsClient.Client)

	privateLinkResourcesClient, err := privatelinkresources.NewPrivateLinkResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrivateLinkResources client: %+v", err)
	}
	configureFunc(privateLinkResourcesClient.Client)

	snapshotsClient, err := snapshots.NewSnapshotsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Snapshots client: %+v", err)
	}
	configureFunc(snapshotsClient.Client)

	volumeGroupsClient, err := volumegroups.NewVolumeGroupsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building VolumeGroups client: %+v", err)
	}
	configureFunc(volumeGroupsClient.Client)

	volumesClient, err := volumes.NewVolumesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Volumes client: %+v", err)
	}
	configureFunc(volumesClient.Client)

	return &Client{
		ElasticSan:                 elasticSanClient,
		ElasticSanSkus:             elasticSanSkusClient,
		ElasticSans:                elasticSansClient,
		PrivateEndpointConnections: privateEndpointConnectionsClient,
		PrivateLinkResources:       privateLinkResourcesClient,
		Snapshots:                  snapshotsClient,
		VolumeGroups:               volumeGroupsClient,
		Volumes:                    volumesClient,
	}, nil
}
