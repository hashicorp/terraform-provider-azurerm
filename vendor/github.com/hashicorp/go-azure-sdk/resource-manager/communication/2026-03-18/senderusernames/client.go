package senderusernames

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SenderUsernamesClient struct {
	Client *resourcemanager.Client
}

func NewSenderUsernamesClientWithBaseURI(sdkApi sdkEnv.Api) (*SenderUsernamesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "senderusernames", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SenderUsernamesClient: %+v", err)
	}

	return &SenderUsernamesClient{
		Client: client,
	}, nil
}
