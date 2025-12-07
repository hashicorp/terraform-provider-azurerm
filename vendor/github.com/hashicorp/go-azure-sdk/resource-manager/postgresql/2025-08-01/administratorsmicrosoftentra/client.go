package administratorsmicrosoftentra

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdministratorsMicrosoftEntraClient struct {
	Client *resourcemanager.Client
}

func NewAdministratorsMicrosoftEntraClientWithBaseURI(sdkApi sdkEnv.Api) (*AdministratorsMicrosoftEntraClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "administratorsmicrosoftentra", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AdministratorsMicrosoftEntraClient: %+v", err)
	}

	return &AdministratorsMicrosoftEntraClient{
		Client: client,
	}, nil
}
