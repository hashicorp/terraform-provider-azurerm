package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2021-09-03-preview/application"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2021-09-03-preview/applicationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2021-09-03-preview/desktop"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2021-09-03-preview/hostpool"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2021-09-03-preview/scalingplan"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2021-09-03-preview/sessionhost"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2021-09-03-preview/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ApplicationGroupsClient *applicationgroup.ApplicationGroupClient
	ApplicationsClient      *application.ApplicationClient
	DesktopsClient          *desktop.DesktopClient
	HostPoolsClient         *hostpool.HostPoolClient
	SessionHostsClient      *sessionhost.SessionHostClient
	ScalingPlansClient      *scalingplan.ScalingPlanClient
	WorkspacesClient        *workspace.WorkspaceClient
}

func NewClient(o *common.ClientOptions) *Client {
	applicationGroupsClient := applicationgroup.NewApplicationGroupClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&applicationGroupsClient.Client, o.ResourceManagerAuthorizer)

	applicationsClient := application.NewApplicationClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&applicationsClient.Client, o.ResourceManagerAuthorizer)

	desktopsClient := desktop.NewDesktopClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&desktopsClient.Client, o.ResourceManagerAuthorizer)

	hostPoolsClient := hostpool.NewHostPoolClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&hostPoolsClient.Client, o.ResourceManagerAuthorizer)

	sessionHostsClient := sessionhost.NewSessionHostClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&sessionHostsClient.Client, o.ResourceManagerAuthorizer)

	scalingPlansClient := scalingplan.NewScalingPlanClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&scalingPlansClient.Client, o.ResourceManagerAuthorizer)

	workspacesClient := workspace.NewWorkspaceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&workspacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ApplicationGroupsClient: &applicationGroupsClient,
		ApplicationsClient:      &applicationsClient,
		DesktopsClient:          &desktopsClient,
		HostPoolsClient:         &hostPoolsClient,
		SessionHostsClient:      &sessionHostsClient,
		ScalingPlansClient:      &scalingPlansClient,
		WorkspacesClient:        &workspacesClient,
	}
}
