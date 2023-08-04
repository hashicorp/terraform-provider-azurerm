// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/assignment"
	"github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/blueprint"
	"github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/publishedblueprint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AssignmentsClient         *assignment.AssignmentClient
	BlueprintsClient          *blueprint.BlueprintClient
	PublishedBlueprintsClient *publishedblueprint.PublishedBlueprintClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	assignmentsClient, err := assignment.NewAssignmentClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Assignment client: %+v", err)
	}
	o.Configure(assignmentsClient.Client, o.Authorizers.ResourceManager)

	blueprintsClient, err := blueprint.NewBlueprintClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Blueprint client: %+v", err)
	}
	o.Configure(blueprintsClient.Client, o.Authorizers.ResourceManager)

	publishedBlueprintsClient, err := publishedblueprint.NewPublishedBlueprintClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building PublishedBlueprint client: %+v", err)
	}
	o.Configure(publishedBlueprintsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AssignmentsClient:         assignmentsClient,
		BlueprintsClient:          blueprintsClient,
		PublishedBlueprintsClient: publishedBlueprintsClient,
	}, nil
}
