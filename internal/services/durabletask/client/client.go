// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/retentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/taskhubs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	SchedulersClient        *schedulers.SchedulersClient
	TaskHubsClient          *taskhubs.TaskHubsClient
	RetentionPoliciesClient *retentionpolicies.RetentionPoliciesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	schedulersClient, err := schedulers.NewSchedulersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Schedulers client: %+v", err)
	}
	o.Configure(schedulersClient.Client, o.Authorizers.ResourceManager)

	taskHubsClient, err := taskhubs.NewTaskHubsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building TaskHubs client: %+v", err)
	}
	o.Configure(taskHubsClient.Client, o.Authorizers.ResourceManager)

	retentionPoliciesClient, err := retentionpolicies.NewRetentionPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building RetentionPolicies client: %+v", err)
	}
	o.Configure(retentionPoliciesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		SchedulersClient:        schedulersClient,
		TaskHubsClient:          taskHubsClient,
		RetentionPoliciesClient: retentionPoliciesClient,
	}, nil
}
