package paloaltonetworkscloudngfws

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PaloAltoNetworksCloudngfwsClient struct {
	Client *resourcemanager.Client
}

func NewPaloAltoNetworksCloudngfwsClientWithBaseURI(sdkApi sdkEnv.Api) (*PaloAltoNetworksCloudngfwsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "paloaltonetworkscloudngfws", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PaloAltoNetworksCloudngfwsClient: %+v", err)
	}

	return &PaloAltoNetworksCloudngfwsClient{
		Client: client,
	}, nil
}
