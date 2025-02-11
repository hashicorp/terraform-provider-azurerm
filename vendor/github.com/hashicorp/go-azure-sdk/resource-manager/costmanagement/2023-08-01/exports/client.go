package exports

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportsClient struct {
	Client *resourcemanager.Client
}

func NewExportsClientWithBaseURI(sdkApi sdkEnv.Api) (*ExportsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "exports", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExportsClient: %+v", err)
	}

	return &ExportsClient{
		Client: client,
	}, nil
}
