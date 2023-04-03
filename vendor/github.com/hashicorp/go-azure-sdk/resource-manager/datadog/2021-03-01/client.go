package v2021_03_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/agreements"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/apikey"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/hosts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/linkedresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/monitoredresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/monitorsresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/refreshsetpasswordlink"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/singlesignon"
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

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	agreementsClient := agreements.NewAgreementsClientWithBaseURI(endpoint)
	configureAuthFunc(&agreementsClient.Client)

	apiKeyClient := apikey.NewApiKeyClientWithBaseURI(endpoint)
	configureAuthFunc(&apiKeyClient.Client)

	hostsClient := hosts.NewHostsClientWithBaseURI(endpoint)
	configureAuthFunc(&hostsClient.Client)

	linkedResourcesClient := linkedresources.NewLinkedResourcesClientWithBaseURI(endpoint)
	configureAuthFunc(&linkedResourcesClient.Client)

	monitoredResourcesClient := monitoredresources.NewMonitoredResourcesClientWithBaseURI(endpoint)
	configureAuthFunc(&monitoredResourcesClient.Client)

	monitorsResourceClient := monitorsresource.NewMonitorsResourceClientWithBaseURI(endpoint)
	configureAuthFunc(&monitorsResourceClient.Client)

	refreshSetPasswordLinkClient := refreshsetpasswordlink.NewRefreshSetPasswordLinkClientWithBaseURI(endpoint)
	configureAuthFunc(&refreshSetPasswordLinkClient.Client)

	rulesClient := rules.NewRulesClientWithBaseURI(endpoint)
	configureAuthFunc(&rulesClient.Client)

	singleSignOnClient := singlesignon.NewSingleSignOnClientWithBaseURI(endpoint)
	configureAuthFunc(&singleSignOnClient.Client)

	return Client{
		Agreements:             &agreementsClient,
		ApiKey:                 &apiKeyClient,
		Hosts:                  &hostsClient,
		LinkedResources:        &linkedResourcesClient,
		MonitoredResources:     &monitoredResourcesClient,
		MonitorsResource:       &monitorsResourceClient,
		RefreshSetPasswordLink: &refreshSetPasswordLinkClient,
		Rules:                  &rulesClient,
		SingleSignOn:           &singleSignOnClient,
	}
}
