package v2020_05_15

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/accesspolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/environments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/eventsources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/referencedatasets"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AccessPolicies    *accesspolicies.AccessPoliciesClient
	Environments      *environments.EnvironmentsClient
	EventSources      *eventsources.EventSourcesClient
	ReferenceDataSets *referencedatasets.ReferenceDataSetsClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	accessPoliciesClient, err := accesspolicies.NewAccessPoliciesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AccessPolicies client: %+v", err)
	}
	configureFunc(accessPoliciesClient.Client)

	environmentsClient, err := environments.NewEnvironmentsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Environments client: %+v", err)
	}
	configureFunc(environmentsClient.Client)

	eventSourcesClient, err := eventsources.NewEventSourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building EventSources client: %+v", err)
	}
	configureFunc(eventSourcesClient.Client)

	referenceDataSetsClient, err := referencedatasets.NewReferenceDataSetsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ReferenceDataSets client: %+v", err)
	}
	configureFunc(referenceDataSetsClient.Client)

	return &Client{
		AccessPolicies:    accessPoliciesClient,
		Environments:      environmentsClient,
		EventSources:      eventSourcesClient,
		ReferenceDataSets: referenceDataSetsClient,
	}, nil
}
