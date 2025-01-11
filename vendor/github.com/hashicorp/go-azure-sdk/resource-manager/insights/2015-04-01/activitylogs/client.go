package activitylogs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActivityLogsClient struct {
	Client *resourcemanager.Client
}

func NewActivityLogsClientWithBaseURI(sdkApi sdkEnv.Api) (*ActivityLogsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "activitylogs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ActivityLogsClient: %+v", err)
	}

	return &ActivityLogsClient{
		Client: client,
	}, nil
}
