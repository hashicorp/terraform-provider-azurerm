package client

import (
	providers "github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2019-06-01-preview/templatespecs"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-09-01/locks"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DeploymentsClient           *resources.DeploymentsClient
	GroupsClient                *resources.GroupsClient
	LocksClient                 *locks.ManagementLocksClient
	ProvidersClient             *providers.ProvidersClient
	ResourcesClient             *resources.Client
	TemplateSpecsVersionsClient *templatespecs.VersionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	deploymentsClient := resources.NewDeploymentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&deploymentsClient.Client, o.ResourceManagerAuthorizer)

	groupsClient := resources.NewGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&groupsClient.Client, o.ResourceManagerAuthorizer)

	locksClient := locks.NewManagementLocksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&locksClient.Client, o.ResourceManagerAuthorizer)

	// this has to come from the Profile since this is shared with Stack
	providersClient := providers.NewProvidersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&providersClient.Client, o.ResourceManagerAuthorizer)

	resourcesClient := resources.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&resourcesClient.Client, o.ResourceManagerAuthorizer)

	templatespecsVersionsClient := templatespecs.NewVersionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&templatespecsVersionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		GroupsClient:                &groupsClient,
		DeploymentsClient:           &deploymentsClient,
		LocksClient:                 &locksClient,
		ProvidersClient:             &providersClient,
		ResourcesClient:             &resourcesClient,
		TemplateSpecsVersionsClient: &templatespecsVersionsClient,
	}
}
