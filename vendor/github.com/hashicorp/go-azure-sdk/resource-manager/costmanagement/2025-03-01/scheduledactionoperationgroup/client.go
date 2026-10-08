package scheduledactionoperationgroup

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduledActionOperationGroupClient struct {
	Client *resourcemanager.Client
}

func NewScheduledActionOperationGroupClientWithBaseURI(sdkApi sdkEnv.Api) (*ScheduledActionOperationGroupClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "scheduledactionoperationgroup", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ScheduledActionOperationGroupClient: %+v", err)
	}

	return &ScheduledActionOperationGroupClient{
		Client: client,
	}, nil
}
