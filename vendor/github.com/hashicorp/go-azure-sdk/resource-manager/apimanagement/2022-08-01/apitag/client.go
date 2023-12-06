package apitag

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiTagClient struct {
	Client *resourcemanager.Client
}

func NewApiTagClientWithBaseURI(sdkApi sdkEnv.Api) (*ApiTagClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "apitag", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ApiTagClient: %+v", err)
	}

	return &ApiTagClient{
		Client: client,
	}, nil
}
