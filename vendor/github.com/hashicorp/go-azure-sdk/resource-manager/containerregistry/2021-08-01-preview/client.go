package v2021_08_01_preview

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/connectedregistries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/exportpipelines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/importpipelines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/operation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/pipelineruns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/replications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/scopemaps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/tokens"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/webhooks"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	ConnectedRegistries        *connectedregistries.ConnectedRegistriesClient
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

func NewClientWithBaseURI(api environments.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	connectedRegistriesClient, err := connectedregistries.NewConnectedRegistriesClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building ConnectedRegistries client: %+v", err)
	}
	configureFunc(connectedRegistriesClient.Client)

	exportPipelinesClient, err := exportpipelines.NewExportPipelinesClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building ExportPipelines client: %+v", err)
	}
	configureFunc(exportPipelinesClient.Client)

	importPipelinesClient, err := importpipelines.NewImportPipelinesClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building ImportPipelines client: %+v", err)
	}
	configureFunc(importPipelinesClient.Client)

	operationClient, err := operation.NewOperationClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building Operation client: %+v", err)
	}
	configureFunc(operationClient.Client)

	pipelineRunsClient, err := pipelineruns.NewPipelineRunsClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building PipelineRuns client: %+v", err)
	}
	configureFunc(pipelineRunsClient.Client)

	privateEndpointConnectionsClient, err := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building PrivateEndpointConnections client: %+v", err)
	}
	configureFunc(privateEndpointConnectionsClient.Client)

	registriesClient, err := registries.NewRegistriesClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building Registries client: %+v", err)
	}
	configureFunc(registriesClient.Client)

	replicationsClient, err := replications.NewReplicationsClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building Replications client: %+v", err)
	}
	configureFunc(replicationsClient.Client)

	scopeMapsClient, err := scopemaps.NewScopeMapsClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building ScopeMaps client: %+v", err)
	}
	configureFunc(scopeMapsClient.Client)

	tokensClient, err := tokens.NewTokensClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building Tokens client: %+v", err)
	}
	configureFunc(tokensClient.Client)

	webHooksClient, err := webhooks.NewWebHooksClientWithBaseURI(api)
	if err != nil {
		return nil, fmt.Errorf("building WebHooks client: %+v", err)
	}
	configureFunc(webHooksClient.Client)

	return &Client{
		ConnectedRegistries:        connectedRegistriesClient,
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
