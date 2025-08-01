package privatedns

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateDNSClient struct {
	Client *resourcemanager.Client
}

func NewPrivateDNSClientWithBaseURI(sdkApi sdkEnv.Api) (*PrivateDNSClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "privatedns", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PrivateDNSClient: %+v", err)
	}

	return &PrivateDNSClient{
		Client: client,
	}, nil
}
