package v2022_12_01

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2022-12-01/arcsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2022-12-01/cluster"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2022-12-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2022-12-01/extensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2022-12-01/offers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2022-12-01/publishers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2022-12-01/skuses"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2022-12-01/updateruns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2022-12-01/updates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2022-12-01/updatesummaries"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	arcSettingsClient := arcsettings.NewArcSettingsClientWithBaseURI(endpoint)
	configureAuthFunc(&arcSettingsClient.Client)

	clusterClient := cluster.NewClusterClientWithBaseURI(endpoint)
	configureAuthFunc(&clusterClient.Client)

	clustersClient := clusters.NewClustersClientWithBaseURI(endpoint)
	configureAuthFunc(&clustersClient.Client)

	extensionsClient := extensions.NewExtensionsClientWithBaseURI(endpoint)
	configureAuthFunc(&extensionsClient.Client)

	offersClient := offers.NewOffersClientWithBaseURI(endpoint)
	configureAuthFunc(&offersClient.Client)

	publishersClient := publishers.NewPublishersClientWithBaseURI(endpoint)
	configureAuthFunc(&publishersClient.Client)

	skusesClient := skuses.NewSkusesClientWithBaseURI(endpoint)
	configureAuthFunc(&skusesClient.Client)

	updateRunsClient := updateruns.NewUpdateRunsClientWithBaseURI(endpoint)
	configureAuthFunc(&updateRunsClient.Client)

	updateSummariesClient := updatesummaries.NewUpdateSummariesClientWithBaseURI(endpoint)
	configureAuthFunc(&updateSummariesClient.Client)

	updatesClient := updates.NewUpdatesClientWithBaseURI(endpoint)
	configureAuthFunc(&updatesClient.Client)

	return Client{
		ArcSettings:     &arcSettingsClient,
		Cluster:         &clusterClient,
		Clusters:        &clustersClient,
		Extensions:      &extensionsClient,
		Offers:          &offersClient,
		Publishers:      &publishersClient,
		Skuses:          &skusesClient,
		UpdateRuns:      &updateRunsClient,
		UpdateSummaries: &updateSummariesClient,
		Updates:         &updatesClient,
	}
}
