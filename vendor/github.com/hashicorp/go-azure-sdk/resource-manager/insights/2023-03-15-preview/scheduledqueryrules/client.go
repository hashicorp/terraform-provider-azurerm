package scheduledqueryrules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduledQueryRulesClient struct {
	Client *resourcemanager.Client
}

func NewScheduledQueryRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*ScheduledQueryRulesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "scheduledqueryrules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ScheduledQueryRulesClient: %+v", err)
	}

	return &ScheduledQueryRulesClient{
		Client: client,
	}, nil
}
