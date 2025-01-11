package packetcorecontrolplane

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PacketCoreControlPlaneClient struct {
	Client *resourcemanager.Client
}

func NewPacketCoreControlPlaneClientWithBaseURI(sdkApi sdkEnv.Api) (*PacketCoreControlPlaneClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "packetcorecontrolplane", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PacketCoreControlPlaneClient: %+v", err)
	}

	return &PacketCoreControlPlaneClient{
		Client: client,
	}, nil
}
