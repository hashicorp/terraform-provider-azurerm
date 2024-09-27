// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/managementlocks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/privatelinkassociation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/resourcemanagementprivatelink"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-10-01/deploymentscripts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-02-01/templatespecversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DeploymentScriptsClient             *deploymentscripts.DeploymentScriptsClient
	FeaturesClient                      *features.FeaturesClient
	LocksClient                         *managementlocks.ManagementLocksClient
	PrivateLinkAssociationClient        *privatelinkassociation.PrivateLinkAssociationClient
	ResourceGroupsClient                *resourcegroups.ResourceGroupsClient
	ResourceManagementPrivateLinkClient *resourcemanagementprivatelink.ResourceManagementPrivateLinkClient
	ResourceProvidersClient             *providers.ProvidersClient
	TemplateSpecsVersionsClient         *templatespecversions.TemplateSpecVersionsClient
	TagsClient                          *tags.TagsClient

	// TODO: these SDK clients use `Azure/azure-sdk-for-go` - we should migrate to `hashicorp/go-azure-sdk`
	// (above) as time allows.
	DeploymentsClient *resources.DeploymentsClient
	ResourcesClient   *resources.Client
	options           *common.ClientOptions

	// Note that the Groups Client which requires additional coordination
	GroupsClient *resources.GroupsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	deploymentScriptsClient, err := deploymentscripts.NewDeploymentScriptsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DeploymentScripts client: %+v", err)
	}
	o.Configure(deploymentScriptsClient.Client, o.Authorizers.ResourceManager)

	featuresClient, err := features.NewFeaturesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Features client: %+v", err)
	}
	o.Configure(featuresClient.Client, o.Authorizers.ResourceManager)

	resourceGroupsClient, err := resourcegroups.NewResourceGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Features client: %+v", err)
	}
	o.Configure(resourceGroupsClient.Client, o.Authorizers.ResourceManager)

	locksClient, err := managementlocks.NewManagementLocksClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ManagementLocks client: %+v", err)
	}
	o.Configure(locksClient.Client, o.Authorizers.ResourceManager)

	privateLinkAssociationClient, err := privatelinkassociation.NewPrivateLinkAssociationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building PrivateLinkAssociation client: %+v", err)
	}
	o.Configure(privateLinkAssociationClient.Client, o.Authorizers.ResourceManager)

	resourceManagementPrivateLinkClient, err := resourcemanagementprivatelink.NewResourceManagementPrivateLinkClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ResourceManagementPrivateLink client: %+v", err)
	}
	o.Configure(resourceManagementPrivateLinkClient.Client, o.Authorizers.ResourceManager)

	resourceProvidersClient, err := providers.NewProvidersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Providers client: %+v", err)
	}
	o.Configure(resourceProvidersClient.Client, o.Authorizers.ResourceManager)

	templateSpecsVersionsClient, err := templatespecversions.NewTemplateSpecVersionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building TemplateSpecVersions client: %+v", err)
	}
	o.Configure(templateSpecsVersionsClient.Client, o.Authorizers.ResourceManager)

	tagsClient, err := tags.NewTagsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Tags client: %+v", err)
	}
	o.Configure(tagsClient.Client, o.Authorizers.ResourceManager)

	// NOTE: these clients use `Azure/azure-sdk-for-go` and can be removed in time
	deploymentsClient := resources.NewDeploymentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&deploymentsClient.Client, o.ResourceManagerAuthorizer)
	groupsClient := resources.NewGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&groupsClient.Client, o.ResourceManagerAuthorizer)
	resourcesClient := resources.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&resourcesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		// These come from `hashicorp/go-azure-sdk`
		DeploymentsClient:                   &deploymentsClient,
		DeploymentScriptsClient:             deploymentScriptsClient,
		FeaturesClient:                      featuresClient,
		LocksClient:                         locksClient,
		PrivateLinkAssociationClient:        privateLinkAssociationClient,
		ResourceManagementPrivateLinkClient: resourceManagementPrivateLinkClient,
		ResourceGroupsClient:                resourceGroupsClient,
		ResourceProvidersClient:             resourceProvidersClient,
		TemplateSpecsVersionsClient:         templateSpecsVersionsClient,
		TagsClient:                          tagsClient,

		// These use `Azure/azure-sdk-for-go`
		GroupsClient:    &groupsClient,
		ResourcesClient: &resourcesClient,
		options:         o,
	}, nil
}
