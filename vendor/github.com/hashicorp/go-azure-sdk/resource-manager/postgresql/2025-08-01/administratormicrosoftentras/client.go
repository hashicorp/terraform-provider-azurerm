package administratormicrosoftentras

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdministratorMicrosoftEntrasClient struct {
	Client *resourcemanager.Client
}

func NewAdministratorMicrosoftEntrasClientWithBaseURI(sdkApi sdkEnv.Api) (*AdministratorMicrosoftEntrasClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "administratormicrosoftentras", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AdministratorMicrosoftEntrasClient: %+v", err)
	}

	return &AdministratorMicrosoftEntrasClient{
		Client: client,
	}, nil
}
