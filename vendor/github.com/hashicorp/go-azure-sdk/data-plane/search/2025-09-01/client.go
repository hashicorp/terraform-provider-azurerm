package v2025_09_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/datasources"
	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/documents"
	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/indexers"
	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/indexes"
	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/service"
	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/skillsets"
	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/synonymmaps"
	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
)

type Client struct {
	DataSources *datasources.DataSourcesClient
	Documents   *documents.DocumentsClient
	Indexers    *indexers.IndexersClient
	Indexes     *indexes.IndexesClient
	Service     *service.ServiceClient
	Skillsets   *skillsets.SkillsetsClient
	SynonymMaps *synonymmaps.SynonymMapsClient
}

func NewClient(configureFunc func(c *dataplane.Client)) (*Client, error) {
	dataSourcesClient, err := datasources.NewDataSourcesClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building DataSources client: %+v", err)
	}
	configureFunc(dataSourcesClient.Client)

	documentsClient, err := documents.NewDocumentsClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Documents client: %+v", err)
	}
	configureFunc(documentsClient.Client)

	indexersClient, err := indexers.NewIndexersClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Indexers client: %+v", err)
	}
	configureFunc(indexersClient.Client)

	indexesClient, err := indexes.NewIndexesClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Indexes client: %+v", err)
	}
	configureFunc(indexesClient.Client)

	serviceClient, err := service.NewServiceClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Service client: %+v", err)
	}
	configureFunc(serviceClient.Client)

	skillsetsClient, err := skillsets.NewSkillsetsClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building Skillsets client: %+v", err)
	}
	configureFunc(skillsetsClient.Client)

	synonymMapsClient, err := synonymmaps.NewSynonymMapsClientUnconfigured()
	if err != nil {
		return nil, fmt.Errorf("building SynonymMaps client: %+v", err)
	}
	configureFunc(synonymMapsClient.Client)

	return &Client{
		DataSources: dataSourcesClient,
		Documents:   documentsClient,
		Indexers:    indexersClient,
		Indexes:     indexesClient,
		Service:     serviceClient,
		Skillsets:   skillsetsClient,
		SynonymMaps: synonymMapsClient,
	}, nil
}
