package notificationrecipientemail

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotificationRecipientEmailClient struct {
	Client *resourcemanager.Client
}

func NewNotificationRecipientEmailClientWithBaseURI(sdkApi sdkEnv.Api) (*NotificationRecipientEmailClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "notificationrecipientemail", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NotificationRecipientEmailClient: %+v", err)
	}

	return &NotificationRecipientEmailClient{
		Client: client,
	}, nil
}
