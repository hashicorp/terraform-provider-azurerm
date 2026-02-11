package v2025_02_15

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/cacertificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/channels"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/clientgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/clients"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/domains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/domaintopics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/eventsubscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/namespacetopics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnerconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnernamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnerregistrations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnertopics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/permissionbindings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/privatelinkresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/systemtopics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/topics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/topicspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/topictypes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/verifiedpartners"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	CaCertificates             *cacertificates.CaCertificatesClient
	Channels                   *channels.ChannelsClient
	ClientGroups               *clientgroups.ClientGroupsClient
	Clients                    *clients.ClientsClient
	DomainTopics               *domaintopics.DomainTopicsClient
	Domains                    *domains.DomainsClient
	EventSubscriptions         *eventsubscriptions.EventSubscriptionsClient
	NamespaceTopics            *namespacetopics.NamespaceTopicsClient
	Namespaces                 *namespaces.NamespacesClient
	PartnerConfigurations      *partnerconfigurations.PartnerConfigurationsClient
	PartnerNamespaces          *partnernamespaces.PartnerNamespacesClient
	PartnerRegistrations       *partnerregistrations.PartnerRegistrationsClient
	PartnerTopics              *partnertopics.PartnerTopicsClient
	PermissionBindings         *permissionbindings.PermissionBindingsClient
	PrivateEndpointConnections *privateendpointconnections.PrivateEndpointConnectionsClient
	PrivateLinkResources       *privatelinkresources.PrivateLinkResourcesClient
	SystemTopics               *systemtopics.SystemTopicsClient
	TopicSpaces                *topicspaces.TopicSpacesClient
	TopicTypes                 *topictypes.TopicTypesClient
	Topics                     *topics.TopicsClient
	VerifiedPartners           *verifiedpartners.VerifiedPartnersClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	caCertificatesClient, err := cacertificates.NewCaCertificatesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CaCertificates client: %+v", err)
	}
	configureFunc(caCertificatesClient.Client)

	channelsClient, err := channels.NewChannelsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Channels client: %+v", err)
	}
	configureFunc(channelsClient.Client)

	clientGroupsClient, err := clientgroups.NewClientGroupsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ClientGroups client: %+v", err)
	}
	configureFunc(clientGroupsClient.Client)

	clientsClient, err := clients.NewClientsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Clients client: %+v", err)
	}
	configureFunc(clientsClient.Client)

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

	namespaceTopicsClient, err := namespacetopics.NewNamespaceTopicsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NamespaceTopics client: %+v", err)
	}
	configureFunc(namespaceTopicsClient.Client)

	namespacesClient, err := namespaces.NewNamespacesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Namespaces client: %+v", err)
	}
	configureFunc(namespacesClient.Client)

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

	permissionBindingsClient, err := permissionbindings.NewPermissionBindingsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PermissionBindings client: %+v", err)
	}
	configureFunc(permissionBindingsClient.Client)

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

	topicSpacesClient, err := topicspaces.NewTopicSpacesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TopicSpaces client: %+v", err)
	}
	configureFunc(topicSpacesClient.Client)

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
		CaCertificates:             caCertificatesClient,
		Channels:                   channelsClient,
		ClientGroups:               clientGroupsClient,
		Clients:                    clientsClient,
		DomainTopics:               domainTopicsClient,
		Domains:                    domainsClient,
		EventSubscriptions:         eventSubscriptionsClient,
		NamespaceTopics:            namespaceTopicsClient,
		Namespaces:                 namespacesClient,
		PartnerConfigurations:      partnerConfigurationsClient,
		PartnerNamespaces:          partnerNamespacesClient,
		PartnerRegistrations:       partnerRegistrationsClient,
		PartnerTopics:              partnerTopicsClient,
		PermissionBindings:         permissionBindingsClient,
		PrivateEndpointConnections: privateEndpointConnectionsClient,
		PrivateLinkResources:       privateLinkResourcesClient,
		SystemTopics:               systemTopicsClient,
		TopicSpaces:                topicSpacesClient,
		TopicTypes:                 topicTypesClient,
		Topics:                     topicsClient,
		VerifiedPartners:           verifiedPartnersClient,
	}, nil
}
