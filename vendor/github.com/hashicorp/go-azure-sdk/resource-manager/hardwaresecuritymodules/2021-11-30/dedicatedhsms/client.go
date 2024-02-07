package dedicatedhsms

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedicatedHsmsClient struct {
	Client *resourcemanager.Client
}

func NewDedicatedHsmsClientWithBaseURI(sdkApi sdkEnv.Api) (*DedicatedHsmsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "dedicatedhsms", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DedicatedHsmsClient: %+v", err)
	}

	return &DedicatedHsmsClient{
		Client: client,
	}, nil
}
