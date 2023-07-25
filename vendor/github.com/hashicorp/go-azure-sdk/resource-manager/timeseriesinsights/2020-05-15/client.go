package v2020_05_15

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/accesspolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/environments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/eventsources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/referencedatasets"
)

type Client struct {
	AccessPolicies    *accesspolicies.AccessPoliciesClient
	Environments      *environments.EnvironmentsClient
	EventSources      *eventsources.EventSourcesClient
	ReferenceDataSets *referencedatasets.ReferenceDataSetsClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	accessPoliciesClient := accesspolicies.NewAccessPoliciesClientWithBaseURI(endpoint)
	configureAuthFunc(&accessPoliciesClient.Client)

	environmentsClient := environments.NewEnvironmentsClientWithBaseURI(endpoint)
	configureAuthFunc(&environmentsClient.Client)

	eventSourcesClient := eventsources.NewEventSourcesClientWithBaseURI(endpoint)
	configureAuthFunc(&eventSourcesClient.Client)

	referenceDataSetsClient := referencedatasets.NewReferenceDataSetsClientWithBaseURI(endpoint)
	configureAuthFunc(&referenceDataSetsClient.Client)

	return Client{
		AccessPolicies:    &accessPoliciesClient,
		Environments:      &environmentsClient,
		EventSources:      &eventSourcesClient,
		ReferenceDataSets: &referenceDataSetsClient,
	}
}
