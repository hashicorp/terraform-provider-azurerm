package eventsubscriptions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventSubscriptionsClient struct {
	Client *resourcemanager.Client
}

func NewEventSubscriptionsClientWithBaseURI(sdkApi sdkEnv.Api) (*EventSubscriptionsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "eventsubscriptions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EventSubscriptionsClient: %+v", err)
	}

	return &EventSubscriptionsClient{
		Client: client,
	}, nil
}
