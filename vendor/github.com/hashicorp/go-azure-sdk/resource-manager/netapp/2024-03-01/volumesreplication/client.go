package volumesreplication

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumesReplicationClient struct {
	Client *resourcemanager.Client
}

func NewVolumesReplicationClientWithBaseURI(sdkApi sdkEnv.Api) (*VolumesReplicationClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "volumesreplication", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VolumesReplicationClient: %+v", err)
	}

	return &VolumesReplicationClient{
		Client: client,
	}, nil
}
