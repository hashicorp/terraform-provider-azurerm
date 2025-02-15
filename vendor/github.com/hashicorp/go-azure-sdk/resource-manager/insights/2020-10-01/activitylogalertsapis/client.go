package activitylogalertsapis

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActivityLogAlertsAPIsClient struct {
	Client *resourcemanager.Client
}

func NewActivityLogAlertsAPIsClientWithBaseURI(sdkApi sdkEnv.Api) (*ActivityLogAlertsAPIsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "activitylogalertsapis", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ActivityLogAlertsAPIsClient: %+v", err)
	}

	return &ActivityLogAlertsAPIsClient{
		Client: client,
	}, nil
}
