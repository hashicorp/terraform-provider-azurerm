// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/lab"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/labplan"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/schedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/user"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	LabClient      *lab.LabClient
	LabPlanClient  *labplan.LabPlanClient
	ScheduleClient *schedule.ScheduleClient
	UserClient     *user.UserClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	labClient, err := lab.NewLabClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(labClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building LabClient client: %+v", err)
	}

	labPlanClient, err := labplan.NewLabPlanClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(labPlanClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building LabPlanClient client: %+v", err)
	}

	scheduleClient, err := schedule.NewScheduleClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(scheduleClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ScheduleClient client: %+v", err)
	}

	userClient, err := user.NewUserClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(userClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building UserClient client: %+v", err)
	}

	return &Client{
		LabClient:      labClient,
		LabPlanClient:  labPlanClient,
		ScheduleClient: scheduleClient,
		UserClient:     userClient,
	}, nil
}
