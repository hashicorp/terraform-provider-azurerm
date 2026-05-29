package v2021_03_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/agreements"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/apikey"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/hosts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/linkedresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/monitoredresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/monitorsresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/refreshsetpasswordlink"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/singlesignon"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Agreements             *agreements.AgreementsClient
	ApiKey                 *apikey.ApiKeyClient
	Hosts                  *hosts.HostsClient
	LinkedResources        *linkedresources.LinkedResourcesClient
	MonitoredResources     *monitoredresources.MonitoredResourcesClient
	MonitorsResource       *monitorsresource.MonitorsResourceClient
	RefreshSetPasswordLink *refreshsetpasswordlink.RefreshSetPasswordLinkClient
	Rules                  *rules.RulesClient
	SingleSignOn           *singlesignon.SingleSignOnClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	agreementsClient, err := agreements.NewAgreementsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Agreements client: %+v", err)
	}
	configureFunc(agreementsClient.Client)

	apiKeyClient, err := apikey.NewApiKeyClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ApiKey client: %+v", err)
	}
	configureFunc(apiKeyClient.Client)

	hostsClient, err := hosts.NewHostsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Hosts client: %+v", err)
	}
	configureFunc(hostsClient.Client)

	linkedResourcesClient, err := linkedresources.NewLinkedResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LinkedResources client: %+v", err)
	}
	configureFunc(linkedResourcesClient.Client)

	monitoredResourcesClient, err := monitoredresources.NewMonitoredResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building MonitoredResources client: %+v", err)
	}
	configureFunc(monitoredResourcesClient.Client)

	monitorsResourceClient, err := monitorsresource.NewMonitorsResourceClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building MonitorsResource client: %+v", err)
	}
	configureFunc(monitorsResourceClient.Client)

	refreshSetPasswordLinkClient, err := refreshsetpasswordlink.NewRefreshSetPasswordLinkClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building RefreshSetPasswordLink client: %+v", err)
	}
	configureFunc(refreshSetPasswordLinkClient.Client)

	rulesClient, err := rules.NewRulesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Rules client: %+v", err)
	}
	configureFunc(rulesClient.Client)

	singleSignOnClient, err := singlesignon.NewSingleSignOnClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SingleSignOn client: %+v", err)
	}
	configureFunc(singleSignOnClient.Client)

	return &Client{
		Agreements:             agreementsClient,
		ApiKey:                 apiKeyClient,
		Hosts:                  hostsClient,
		LinkedResources:        linkedResourcesClient,
		MonitoredResources:     monitoredResourcesClient,
		MonitorsResource:       monitorsResourceClient,
		RefreshSetPasswordLink: refreshSetPasswordLinkClient,
		Rules:                  rulesClient,
		SingleSignOn:           singleSignOnClient,
	}, nil
}
