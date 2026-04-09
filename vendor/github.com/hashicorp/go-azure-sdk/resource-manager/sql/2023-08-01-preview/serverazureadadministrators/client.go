package serverazureadadministrators

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerAzureADAdministratorsClient struct {
	Client *resourcemanager.Client
}

func NewServerAzureADAdministratorsClientWithBaseURI(sdkApi sdkEnv.Api) (*ServerAzureADAdministratorsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "serverazureadadministrators", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServerAzureADAdministratorsClient: %+v", err)
	}

	return &ServerAzureADAdministratorsClient{
		Client: client,
	}, nil
}
