// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2019-06-01-preview/templatespecs" // nolint: staticcheck
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"                     // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/managementlocks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-10-01/deploymentscripts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DeploymentsClient           *resources.DeploymentsClient
	DeploymentScriptsClient     *deploymentscripts.DeploymentScriptsClient
	FeaturesClient              *features.FeaturesClient
	GroupsClient                *resources.GroupsClient
	LocksClient                 *managementlocks.ManagementLocksClient
	ResourceProvidersClient     *providers.ProvidersClient
	ResourcesClient             *resources.Client
	TagsClient                  *resources.TagsClient
	TemplateSpecsVersionsClient *templatespecs.VersionsClient

	options *common.ClientOptions
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	deploymentsClient := resources.NewDeploymentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&deploymentsClient.Client, o.ResourceManagerAuthorizer)

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

	groupsClient := resources.NewGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&groupsClient.Client, o.ResourceManagerAuthorizer)

	locksClient, err := managementlocks.NewManagementLocksClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ManagementLocks client: %+v", err)
	}
	o.Configure(locksClient.Client, o.Authorizers.ResourceManager)

	resourceProvidersClient, err := providers.NewProvidersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Providers client: %+v", err)
	}
	o.Configure(resourceProvidersClient.Client, o.Authorizers.ResourceManager)

	resourcesClient := resources.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&resourcesClient.Client, o.ResourceManagerAuthorizer)

	templatespecsVersionsClient := templatespecs.NewVersionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&templatespecsVersionsClient.Client, o.ResourceManagerAuthorizer)

	tagsClient := resources.NewTagsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tagsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		GroupsClient:                &groupsClient,
		DeploymentsClient:           &deploymentsClient,
		DeploymentScriptsClient:     deploymentScriptsClient,
		FeaturesClient:              featuresClient,
		LocksClient:                 locksClient,
		ResourceProvidersClient:     resourceProvidersClient,
		ResourcesClient:             &resourcesClient,
		TagsClient:                  &tagsClient,
		TemplateSpecsVersionsClient: &templatespecsVersionsClient,

		options: o,
	}, nil
}

func (c Client) TagsClientForSubscription(subscriptionID string) *resources.TagsClient {
	// TODO: this method can be removed once this is moved to using `hashicorp/go-azure-sdk`
	tagsClient := resources.NewTagsClientWithBaseURI(c.options.ResourceManagerEndpoint, subscriptionID)
	c.options.ConfigureClient(&tagsClient.Client, c.options.ResourceManagerAuthorizer)
	return &tagsClient
}
