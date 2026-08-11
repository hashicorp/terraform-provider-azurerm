package scalingplan

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScalingPlanClient struct {
	Client *resourcemanager.Client
}

func NewScalingPlanClientWithBaseURI(sdkApi sdkEnv.Api) (*ScalingPlanClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "scalingplan", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ScalingPlanClient: %+v", err)
	}

	return &ScalingPlanClient{
		Client: client,
	}, nil
}
