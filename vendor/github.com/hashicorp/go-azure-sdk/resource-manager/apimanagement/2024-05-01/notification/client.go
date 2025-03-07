package notification

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotificationClient struct {
	Client *resourcemanager.Client
}

func NewNotificationClientWithBaseURI(sdkApi sdkEnv.Api) (*NotificationClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "notification", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NotificationClient: %+v", err)
	}

	return &NotificationClient{
		Client: client,
	}, nil
}
