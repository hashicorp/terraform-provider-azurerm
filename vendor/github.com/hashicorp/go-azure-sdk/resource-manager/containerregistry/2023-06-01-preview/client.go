package v2023_06_01_preview

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/archives"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/archiveversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/cacherules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/connectedregistries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/credentialsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/exportpipelines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/importpipelines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/operation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/pipelineruns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/replications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/scopemaps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/tokens"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/webhooks"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	ArchiveVersions            *archiveversions.ArchiveVersionsClient
	Archives                   *archives.ArchivesClient
	CacheRules                 *cacherules.CacheRulesClient
	ConnectedRegistries        *connectedregistries.ConnectedRegistriesClient
	CredentialSets             *credentialsets.CredentialSetsClient
	ExportPipelines            *exportpipelines.ExportPipelinesClient
	ImportPipelines            *importpipelines.ImportPipelinesClient
	Operation                  *operation.OperationClient
	PipelineRuns               *pipelineruns.PipelineRunsClient
	PrivateEndpointConnections *privateendpointconnections.PrivateEndpointConnectionsClient
	Registries                 *registries.RegistriesClient
	Replications               *replications.ReplicationsClient
	ScopeMaps                  *scopemaps.ScopeMapsClient
	Tokens                     *tokens.TokensClient
	WebHooks                   *webhooks.WebHooksClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	archiveVersionsClient, err := archiveversions.NewArchiveVersionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ArchiveVersions client: %+v", err)
	}
	configureFunc(archiveVersionsClient.Client)

	archivesClient, err := archives.NewArchivesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Archives client: %+v", err)
	}
	configureFunc(archivesClient.Client)

	cacheRulesClient, err := cacherules.NewCacheRulesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CacheRules client: %+v", err)
	}
	configureFunc(cacheRulesClient.Client)

	connectedRegistriesClient, err := connectedregistries.NewConnectedRegistriesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ConnectedRegistries client: %+v", err)
	}
	configureFunc(connectedRegistriesClient.Client)

	credentialSetsClient, err := credentialsets.NewCredentialSetsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CredentialSets client: %+v", err)
	}
	configureFunc(credentialSetsClient.Client)

	exportPipelinesClient, err := exportpipelines.NewExportPipelinesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ExportPipelines client: %+v", err)
	}
	configureFunc(exportPipelinesClient.Client)

	importPipelinesClient, err := importpipelines.NewImportPipelinesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ImportPipelines client: %+v", err)
	}
	configureFunc(importPipelinesClient.Client)

	operationClient, err := operation.NewOperationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Operation client: %+v", err)
	}
	configureFunc(operationClient.Client)

	pipelineRunsClient, err := pipelineruns.NewPipelineRunsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PipelineRuns client: %+v", err)
	}
	configureFunc(pipelineRunsClient.Client)

	privateEndpointConnectionsClient, err := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrivateEndpointConnections client: %+v", err)
	}
	configureFunc(privateEndpointConnectionsClient.Client)

	registriesClient, err := registries.NewRegistriesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Registries client: %+v", err)
	}
	configureFunc(registriesClient.Client)

	replicationsClient, err := replications.NewReplicationsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Replications client: %+v", err)
	}
	configureFunc(replicationsClient.Client)

	scopeMapsClient, err := scopemaps.NewScopeMapsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ScopeMaps client: %+v", err)
	}
	configureFunc(scopeMapsClient.Client)

	tokensClient, err := tokens.NewTokensClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Tokens client: %+v", err)
	}
	configureFunc(tokensClient.Client)

	webHooksClient, err := webhooks.NewWebHooksClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building WebHooks client: %+v", err)
	}
	configureFunc(webHooksClient.Client)

	return &Client{
		ArchiveVersions:            archiveVersionsClient,
		Archives:                   archivesClient,
		CacheRules:                 cacheRulesClient,
		ConnectedRegistries:        connectedRegistriesClient,
		CredentialSets:             credentialSetsClient,
		ExportPipelines:            exportPipelinesClient,
		ImportPipelines:            importPipelinesClient,
		Operation:                  operationClient,
		PipelineRuns:               pipelineRunsClient,
		PrivateEndpointConnections: privateEndpointConnectionsClient,
		Registries:                 registriesClient,
		Replications:               replicationsClient,
		ScopeMaps:                  scopeMapsClient,
		Tokens:                     tokensClient,
		WebHooks:                   webHooksClient,
	}, nil
}
