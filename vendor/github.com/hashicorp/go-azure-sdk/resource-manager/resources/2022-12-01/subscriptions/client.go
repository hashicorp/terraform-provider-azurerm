package subscriptions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionsClient struct {
	Client *resourcemanager.Client
}

func NewSubscriptionsClientWithBaseURI(sdkApi sdkEnv.Api) (*SubscriptionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "subscriptions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SubscriptionsClient: %+v", err)
	}

	return &SubscriptionsClient{
		Client: client,
	}, nil
}
