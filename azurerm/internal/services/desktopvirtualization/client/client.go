package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/desktopvirtualization/mgmt/2020-11-02-preview/desktopvirtualization"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ApplicationGroupsClient *desktopvirtualization.ApplicationGroupsClient
	ApplicationsClient      *desktopvirtualization.ApplicationsClient
	DesktopsClient          *desktopvirtualization.DesktopsClient
	HostPoolsClient         *desktopvirtualization.HostPoolsClient
	OperationsClient        *desktopvirtualization.OperationsClient
	SessionHostsClient      *desktopvirtualization.SessionHostsClient
	WorkspacesClient        *desktopvirtualization.WorkspacesClient
}

// NewClient - New client for desktop virtualization
func NewClient(o *common.ClientOptions) *Client {
	ApplicationGroupsClient := desktopvirtualization.NewApplicationGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ApplicationGroupsClient.Client, o.ResourceManagerAuthorizer)

	ApplicationsClient := desktopvirtualization.NewApplicationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ApplicationsClient.Client, o.ResourceManagerAuthorizer)

	DesktopsClient := desktopvirtualization.NewDesktopsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DesktopsClient.Client, o.ResourceManagerAuthorizer)

	HostPoolsClient := desktopvirtualization.NewHostPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HostPoolsClient.Client, o.ResourceManagerAuthorizer)

	OperationsClient := desktopvirtualization.NewOperationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&OperationsClient.Client, o.ResourceManagerAuthorizer)

	SessionHostsClient := desktopvirtualization.NewSessionHostsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SessionHostsClient.Client, o.ResourceManagerAuthorizer)

	WorkspacesClient := desktopvirtualization.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ApplicationGroupsClient: &ApplicationGroupsClient,
		ApplicationsClient:      &ApplicationsClient,
		DesktopsClient:          &DesktopsClient,
		HostPoolsClient:         &HostPoolsClient,
		OperationsClient:        &OperationsClient,
		SessionHostsClient:      &SessionHostsClient,
		WorkspacesClient:        &WorkspacesClient,
	}
}
