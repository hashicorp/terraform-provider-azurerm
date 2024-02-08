package blueprint

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlueprintClient struct {
	Client *resourcemanager.Client
}

func NewBlueprintClientWithBaseURI(sdkApi sdkEnv.Api) (*BlueprintClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "blueprint", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BlueprintClient: %+v", err)
	}

	return &BlueprintClient{
		Client: client,
	}, nil
}
