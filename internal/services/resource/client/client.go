// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	azureResources "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/managementlocks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/privatelinkassociation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/resourcemanagementprivatelink"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-10-01/deploymentscripts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-02-01/templatespecversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/deployments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2024-03-01/deploymentstacksatmanagementgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2024-03-01/deploymentstacksatresourcegroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2024-03-01/deploymentstacksatsubscription"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DeploymentsClient                     *deployments.DeploymentsClient
	DeploymentScriptsClient               *deploymentscripts.DeploymentScriptsClient
	DeploymentStacksManagementGroupClient *deploymentstacksatmanagementgroup.DeploymentStacksAtManagementGroupClient
	DeploymentStacksResourceGroupClient   *deploymentstacksatresourcegroup.DeploymentStacksAtResourceGroupClient
	DeploymentStacksSubscriptionClient    *deploymentstacksatsubscription.DeploymentStacksAtSubscriptionClient
	FeaturesClient                        *features.FeaturesClient
	LocksClient                           *managementlocks.ManagementLocksClient
	PrivateLinkAssociationClient          *privatelinkassociation.PrivateLinkAssociationClient
	ResourcesClient                       *resources.ResourcesClient
	ResourceGroupsClient                  *resourcegroups.ResourceGroupsClient
	ResourceManagementPrivateLinkClient   *resourcemanagementprivatelink.ResourceManagementPrivateLinkClient
	ResourceProvidersClient               *providers.ProvidersClient
	TemplateSpecsVersionsClient           *templatespecversions.TemplateSpecVersionsClient
	TagsClient                            *tags.TagsClient

	// TODO: these SDK clients use `Azure/azure-sdk-for-go` - we should migrate to `hashicorp/go-azure-sdk`
	// (above) as time allows.
	options                 *common.ClientOptions
	LegacyResourcesClient   *azureResources.Client
	LegacyDeploymentsClient *azureResources.DeploymentsClient

	// Note that the Groups Client which requires additional coordination
	GroupsClient *azureResources.GroupsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	deploymentsClient, err := deployments.NewDeploymentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Deployments client: %+v", err)
	}
	o.Configure(deploymentsClient.Client, o.Authorizers.ResourceManager)

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

	resourcesClient, err := resources.NewResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Resource client: %+v", err)
	}
	o.Configure(resourcesClient.Client, o.Authorizers.ResourceManager)

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

	deploymentStacksResourceGroupClient, err := deploymentstacksatresourcegroup.NewDeploymentStacksAtResourceGroupClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DeploymentStacksAtResourceGroup client: %+v", err)
	}
	o.Configure(deploymentStacksResourceGroupClient.Client, o.Authorizers.ResourceManager)

	deploymentStacksSubscriptionClient, err := deploymentstacksatsubscription.NewDeploymentStacksAtSubscriptionClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DeploymentStacksAtSubscription client: %+v", err)
	}
	o.Configure(deploymentStacksSubscriptionClient.Client, o.Authorizers.ResourceManager)

	deploymentStacksManagementGroupClient, err := deploymentstacksatmanagementgroup.NewDeploymentStacksAtManagementGroupClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DeploymentStacksAtManagementGroup client: %+v", err)
	}
	o.Configure(deploymentStacksManagementGroupClient.Client, o.Authorizers.ResourceManager)

	// NOTE: This client uses `Azure/azure-sdk-for-go` and can be removed in time
	legacyDeploymentsClient := azureResources.NewDeploymentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&legacyDeploymentsClient.Client, o.ResourceManagerAuthorizer)

	legacyResourcesClient := azureResources.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&legacyResourcesClient.Client, o.ResourceManagerAuthorizer)

	groupsClient := azureResources.NewGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&groupsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		// These come from `hashicorp/go-azure-sdk`
		DeploymentsClient:                     deploymentsClient,
		DeploymentScriptsClient:               deploymentScriptsClient,
		DeploymentStacksManagementGroupClient: deploymentStacksManagementGroupClient,
		DeploymentStacksResourceGroupClient:   deploymentStacksResourceGroupClient,
		DeploymentStacksSubscriptionClient:    deploymentStacksSubscriptionClient,
		FeaturesClient:                        featuresClient,
		LocksClient:                           locksClient,
		PrivateLinkAssociationClient:          privateLinkAssociationClient,
		ResourcesClient:                       resourcesClient,
		ResourceManagementPrivateLinkClient:   resourceManagementPrivateLinkClient,
		ResourceGroupsClient:                  resourceGroupsClient,
		ResourceProvidersClient:               resourceProvidersClient,
		TemplateSpecsVersionsClient:           templateSpecsVersionsClient,
		TagsClient:                            tagsClient,

		// These use `Azure/azure-sdk-for-go`
		LegacyDeploymentsClient: &legacyDeploymentsClient,
		LegacyResourcesClient:   &legacyResourcesClient,
		GroupsClient:            &groupsClient,
		options:                 o,
	}, nil
}
