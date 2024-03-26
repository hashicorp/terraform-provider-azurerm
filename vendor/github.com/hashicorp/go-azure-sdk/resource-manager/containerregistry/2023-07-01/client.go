package v2023_07_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/cacherules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/credentialsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/operation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/replications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/scopemaps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/tokens"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/webhooks"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	CacheRules                 *cacherules.CacheRulesClient
	CredentialSets             *credentialsets.CredentialSetsClient
	Operation                  *operation.OperationClient
	PrivateEndpointConnections *privateendpointconnections.PrivateEndpointConnectionsClient
	Registries                 *registries.RegistriesClient
	Replications               *replications.ReplicationsClient
	ScopeMaps                  *scopemaps.ScopeMapsClient
	Tokens                     *tokens.TokensClient
	WebHooks                   *webhooks.WebHooksClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	cacheRulesClient, err := cacherules.NewCacheRulesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CacheRules client: %+v", err)
	}
	configureFunc(cacheRulesClient.Client)

	credentialSetsClient, err := credentialsets.NewCredentialSetsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CredentialSets client: %+v", err)
	}
	configureFunc(credentialSetsClient.Client)

	operationClient, err := operation.NewOperationClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Operation client: %+v", err)
	}
	configureFunc(operationClient.Client)

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
		CacheRules:                 cacheRulesClient,
		CredentialSets:             credentialSetsClient,
		Operation:                  operationClient,
		PrivateEndpointConnections: privateEndpointConnectionsClient,
		Registries:                 registriesClient,
		Replications:               replicationsClient,
		ScopeMaps:                  scopeMapsClient,
		Tokens:                     tokensClient,
		WebHooks:                   webHooksClient,
	}, nil
}
