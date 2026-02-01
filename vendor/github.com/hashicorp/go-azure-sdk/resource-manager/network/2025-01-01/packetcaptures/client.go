package packetcaptures

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PacketCapturesClient struct {
	Client *resourcemanager.Client
}

func NewPacketCapturesClientWithBaseURI(sdkApi sdkEnv.Api) (*PacketCapturesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "packetcaptures", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PacketCapturesClient: %+v", err)
	}

	return &PacketCapturesClient{
		Client: client,
	}, nil
}
