package v2022_06_15

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/channels"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/domains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/domaintopics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/eventsubscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnerconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnernamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnerregistrations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnertopics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/privatelinkresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/systemtopics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/topics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/topictypes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/verifiedpartners"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	Channels                   *channels.ChannelsClient
	DomainTopics               *domaintopics.DomainTopicsClient
	Domains                    *domains.DomainsClient
	EventSubscriptions         *eventsubscriptions.EventSubscriptionsClient
	PartnerConfigurations      *partnerconfigurations.PartnerConfigurationsClient
	PartnerNamespaces          *partnernamespaces.PartnerNamespacesClient
	PartnerRegistrations       *partnerregistrations.PartnerRegistrationsClient
	PartnerTopics              *partnertopics.PartnerTopicsClient
	PrivateEndpointConnections *privateendpointconnections.PrivateEndpointConnectionsClient
	PrivateLinkResources       *privatelinkresources.PrivateLinkResourcesClient
	SystemTopics               *systemtopics.SystemTopicsClient
	TopicTypes                 *topictypes.TopicTypesClient
	Topics                     *topics.TopicsClient
	VerifiedPartners           *verifiedpartners.VerifiedPartnersClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	channelsClient, err := channels.NewChannelsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Channels client: %+v", err)
	}
	configureFunc(channelsClient.Client)

	domainTopicsClient, err := domaintopics.NewDomainTopicsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DomainTopics client: %+v", err)
	}
	configureFunc(domainTopicsClient.Client)

	domainsClient, err := domains.NewDomainsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Domains client: %+v", err)
	}
	configureFunc(domainsClient.Client)

	eventSubscriptionsClient, err := eventsubscriptions.NewEventSubscriptionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building EventSubscriptions client: %+v", err)
	}
	configureFunc(eventSubscriptionsClient.Client)

	partnerConfigurationsClient, err := partnerconfigurations.NewPartnerConfigurationsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PartnerConfigurations client: %+v", err)
	}
	configureFunc(partnerConfigurationsClient.Client)

	partnerNamespacesClient, err := partnernamespaces.NewPartnerNamespacesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PartnerNamespaces client: %+v", err)
	}
	configureFunc(partnerNamespacesClient.Client)

	partnerRegistrationsClient, err := partnerregistrations.NewPartnerRegistrationsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PartnerRegistrations client: %+v", err)
	}
	configureFunc(partnerRegistrationsClient.Client)

	partnerTopicsClient, err := partnertopics.NewPartnerTopicsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PartnerTopics client: %+v", err)
	}
	configureFunc(partnerTopicsClient.Client)

	privateEndpointConnectionsClient, err := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrivateEndpointConnections client: %+v", err)
	}
	configureFunc(privateEndpointConnectionsClient.Client)

	privateLinkResourcesClient, err := privatelinkresources.NewPrivateLinkResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrivateLinkResources client: %+v", err)
	}
	configureFunc(privateLinkResourcesClient.Client)

	systemTopicsClient, err := systemtopics.NewSystemTopicsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SystemTopics client: %+v", err)
	}
	configureFunc(systemTopicsClient.Client)

	topicTypesClient, err := topictypes.NewTopicTypesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TopicTypes client: %+v", err)
	}
	configureFunc(topicTypesClient.Client)

	topicsClient, err := topics.NewTopicsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Topics client: %+v", err)
	}
	configureFunc(topicsClient.Client)

	verifiedPartnersClient, err := verifiedpartners.NewVerifiedPartnersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building VerifiedPartners client: %+v", err)
	}
	configureFunc(verifiedPartnersClient.Client)

	return &Client{
		Channels:                   channelsClient,
		DomainTopics:               domainTopicsClient,
		Domains:                    domainsClient,
		EventSubscriptions:         eventSubscriptionsClient,
		PartnerConfigurations:      partnerConfigurationsClient,
		PartnerNamespaces:          partnerNamespacesClient,
		PartnerRegistrations:       partnerRegistrationsClient,
		PartnerTopics:              partnerTopicsClient,
		PrivateEndpointConnections: privateEndpointConnectionsClient,
		PrivateLinkResources:       privateLinkResourcesClient,
		SystemTopics:               systemTopicsClient,
		TopicTypes:                 topicTypesClient,
		Topics:                     topicsClient,
		VerifiedPartners:           verifiedPartnersClient,
	}, nil
}
