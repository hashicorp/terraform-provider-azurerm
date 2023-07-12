// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/application"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/applicationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/desktop"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/hostpool"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/scalingplan"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/sessionhost"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2022-02-10-preview/workspace"
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

func NewClient(o *common.ClientOptions) (*Client, error) {
	applicationGroupsClient, err := applicationgroup.NewApplicationGroupClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ApplicationGroups Client: %+v", err)
	}
	o.Configure(applicationGroupsClient.Client, o.Authorizers.ResourceManager)

	applicationsClient, err := application.NewApplicationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Applications Client: %+v", err)
	}
	o.Configure(applicationsClient.Client, o.Authorizers.ResourceManager)

	desktopsClient, err := desktop.NewDesktopClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Desktops Client: %+v", err)
	}
	o.Configure(desktopsClient.Client, o.Authorizers.ResourceManager)

	hostPoolsClient, err := hostpool.NewHostPoolClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building HostPools Client: %+v", err)
	}
	o.Configure(hostPoolsClient.Client, o.Authorizers.ResourceManager)

	sessionHostsClient, err := sessionhost.NewSessionHostClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SessionHost Client: %+v", err)
	}
	o.Configure(sessionHostsClient.Client, o.Authorizers.ResourceManager)

	scalingPlansClient, err := scalingplan.NewScalingPlanClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ScalingPlan Client: %+v", err)
	}
	o.Configure(scalingPlansClient.Client, o.Authorizers.ResourceManager)

	workspacesClient, err := workspace.NewWorkspaceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Workspaces Client: %+v", err)
	}
	o.Configure(workspacesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ApplicationGroupsClient: applicationGroupsClient,
		ApplicationsClient:      applicationsClient,
		DesktopsClient:          desktopsClient,
		HostPoolsClient:         hostPoolsClient,
		SessionHostsClient:      sessionHostsClient,
		ScalingPlansClient:      scalingPlansClient,
		WorkspacesClient:        workspacesClient,
	}, nil
}
