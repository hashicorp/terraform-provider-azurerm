package volumes

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumesClient struct {
	Client *resourcemanager.Client
}

func NewVolumesClientWithBaseURI(sdkApi sdkEnv.Api) (*VolumesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "volumes", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VolumesClient: %+v", err)
	}

	return &VolumesClient{
		Client: client,
	}, nil
}
