package v2025_06_11

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogmonitorresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/datadogsinglesignonresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/monitoredsubscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/tagrules"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	DatadogMonitorResources      *datadogmonitorresources.DatadogMonitorResourcesClient
	DatadogSingleSignOnResources *datadogsinglesignonresources.DatadogSingleSignOnResourcesClient
	Datadogs                     *datadogs.DatadogsClient
	MonitoredSubscriptions       *monitoredsubscriptions.MonitoredSubscriptionsClient
	TagRules                     *tagrules.TagRulesClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	datadogMonitorResourcesClient, err := datadogmonitorresources.NewDatadogMonitorResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DatadogMonitorResources client: %+v", err)
	}
	configureFunc(datadogMonitorResourcesClient.Client)

	datadogSingleSignOnResourcesClient, err := datadogsinglesignonresources.NewDatadogSingleSignOnResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DatadogSingleSignOnResources client: %+v", err)
	}
	configureFunc(datadogSingleSignOnResourcesClient.Client)

	datadogsClient, err := datadogs.NewDatadogsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Datadogs client: %+v", err)
	}
	configureFunc(datadogsClient.Client)

	monitoredSubscriptionsClient, err := monitoredsubscriptions.NewMonitoredSubscriptionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building MonitoredSubscriptions client: %+v", err)
	}
	configureFunc(monitoredSubscriptionsClient.Client)

	tagRulesClient, err := tagrules.NewTagRulesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TagRules client: %+v", err)
	}
	configureFunc(tagRulesClient.Client)

	return &Client{
		DatadogMonitorResources:      datadogMonitorResourcesClient,
		DatadogSingleSignOnResources: datadogSingleSignOnResourcesClient,
		Datadogs:                     datadogsClient,
		MonitoredSubscriptions:       monitoredSubscriptionsClient,
		TagRules:                     tagRulesClient,
	}, nil
}
