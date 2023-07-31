package v2023_03_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/arcsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/cluster"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/extensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/offers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/publishers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/skuses"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/updateruns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/updates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/updatesummaries"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	ArcSettings     *arcsettings.ArcSettingsClient
	Cluster         *cluster.ClusterClient
	Clusters        *clusters.ClustersClient
	Extensions      *extensions.ExtensionsClient
	Offers          *offers.OffersClient
	Publishers      *publishers.PublishersClient
	Skuses          *skuses.SkusesClient
	UpdateRuns      *updateruns.UpdateRunsClient
	UpdateSummaries *updatesummaries.UpdateSummariesClient
	Updates         *updates.UpdatesClient
}

func NewClientWithBaseURI(api environments.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	arcSettingsClient, err := arcsettings.NewArcSettingsClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building ArcSettings client: %+v", err)
	}
	configureFunc(arcSettingsClient.Client)

	clusterClient, err := cluster.NewClusterClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building Cluster client: %+v", err)
	}
	configureFunc(clusterClient.Client)

	clustersClient, err := clusters.NewClustersClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building Clusters client: %+v", err)
	}
	configureFunc(clustersClient.Client)

	extensionsClient, err := extensions.NewExtensionsClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building Extensions client: %+v", err)
	}
	configureFunc(extensionsClient.Client)

	offersClient, err := offers.NewOffersClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building Offers client: %+v", err)
	}
	configureFunc(offersClient.Client)

	publishersClient, err := publishers.NewPublishersClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building Publishers client: %+v", err)
	}
	configureFunc(publishersClient.Client)

	skusesClient, err := skuses.NewSkusesClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building Skuses client: %+v", err)
	}
	configureFunc(skusesClient.Client)

	updateRunsClient, err := updateruns.NewUpdateRunsClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building UpdateRuns client: %+v", err)
	}
	configureFunc(updateRunsClient.Client)

	updateSummariesClient, err := updatesummaries.NewUpdateSummariesClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building UpdateSummaries client: %+v", err)
	}
	configureFunc(updateSummariesClient.Client)

	updatesClient, err := updates.NewUpdatesClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building Updates client: %+v", err)
	}
	configureFunc(updatesClient.Client)

	return &Client{
		ArcSettings:     arcSettingsClient,
		Cluster:         clusterClient,
		Clusters:        clustersClient,
		Extensions:      extensionsClient,
		Offers:          offersClient,
		Publishers:      publishersClient,
		Skuses:          skusesClient,
		UpdateRuns:      updateRunsClient,
		UpdateSummaries: updateSummariesClient,
		Updates:         updatesClient,
	}, nil
}
