package client

import (
	providers "github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2019-06-01-preview/templatespecs"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2015-12-01/features"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-09-01/locks"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DeploymentsClient           *resources.DeploymentsClient
	FeaturesClient              *features.Client
	GroupsClient                *resources.GroupsClient
	LocksClient                 *locks.ManagementLocksClient
	ProvidersClient             *providers.ProvidersClient
	ResourceProvidersClient     *resources.ProvidersClient
	ResourcesClient             *resources.Client
	TagsClient                  *resources.TagsClient
	TemplateSpecsVersionsClient *templatespecs.VersionsClient

	options *common.ClientOptions
}

func NewClient(o *common.ClientOptions) *Client {
	deploymentsClient := resources.NewDeploymentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&deploymentsClient.Client, o.ResourceManagerAuthorizer)

	featuresClient := features.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&featuresClient.Client, o.ResourceManagerAuthorizer)

	groupsClient := resources.NewGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&groupsClient.Client, o.ResourceManagerAuthorizer)

	locksClient := locks.NewManagementLocksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&locksClient.Client, o.ResourceManagerAuthorizer)

	// this has to come from the Profile since this is shared with Stack
	providersClient := providers.NewProvidersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&providersClient.Client, o.ResourceManagerAuthorizer)

	// add a secondary ProvidersClient to use latest resources sdk
	resourceProvidersClient := resources.NewProvidersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&resourceProvidersClient.Client, o.ResourceManagerAuthorizer)

	resourcesClient := resources.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&resourcesClient.Client, o.ResourceManagerAuthorizer)

	templatespecsVersionsClient := templatespecs.NewVersionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&templatespecsVersionsClient.Client, o.ResourceManagerAuthorizer)

	tagsClient := resources.NewTagsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tagsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		GroupsClient:                &groupsClient,
		DeploymentsClient:           &deploymentsClient,
		FeaturesClient:              &featuresClient,
		LocksClient:                 &locksClient,
		ProvidersClient:             &providersClient,
		ResourceProvidersClient:     &resourceProvidersClient,
		ResourcesClient:             &resourcesClient,
		TagsClient:                  &tagsClient,
		TemplateSpecsVersionsClient: &templatespecsVersionsClient,

		options: o,
	}
}

func (c Client) TagsClientForSubscription(subscriptionID string) *resources.TagsClient {
	tagsClient := resources.NewTagsClientWithBaseURI(c.options.ResourceManagerEndpoint, subscriptionID)
	c.options.ConfigureClient(&tagsClient.Client, c.options.ResourceManagerAuthorizer)
	return &tagsClient
}
