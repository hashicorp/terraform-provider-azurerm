package usertoken

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserTokenClient struct {
	Client *resourcemanager.Client
}

func NewUserTokenClientWithBaseURI(sdkApi sdkEnv.Api) (*UserTokenClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "usertoken", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating UserTokenClient: %+v", err)
	}

	return &UserTokenClient{
		Client: client,
	}, nil
}
