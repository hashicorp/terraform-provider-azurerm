package administrators

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdministratorsClient struct {
	Client *resourcemanager.Client
}

func NewAdministratorsClientWithBaseURI(sdkApi sdkEnv.Api) (*AdministratorsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "administrators", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AdministratorsClient: %+v", err)
	}

	return &AdministratorsClient{
		Client: client,
	}, nil
}
