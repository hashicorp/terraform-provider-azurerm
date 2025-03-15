package userconfirmationpasswordsend

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserConfirmationPasswordSendClient struct {
	Client *resourcemanager.Client
}

func NewUserConfirmationPasswordSendClientWithBaseURI(sdkApi sdkEnv.Api) (*UserConfirmationPasswordSendClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "userconfirmationpasswordsend", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating UserConfirmationPasswordSendClient: %+v", err)
	}

	return &UserConfirmationPasswordSendClient{
		Client: client,
	}, nil
}
