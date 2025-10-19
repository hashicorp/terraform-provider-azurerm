package deploymentstacksatsubscription

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentStacksAtSubscriptionClient struct {
	Client *resourcemanager.Client
}

func NewDeploymentStacksAtSubscriptionClientWithBaseURI(sdkApi sdkEnv.Api) (*DeploymentStacksAtSubscriptionClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "deploymentstacksatsubscription", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DeploymentStacksAtSubscriptionClient: %+v", err)
	}

	return &DeploymentStacksAtSubscriptionClient{
		Client: client,
	}, nil
}
