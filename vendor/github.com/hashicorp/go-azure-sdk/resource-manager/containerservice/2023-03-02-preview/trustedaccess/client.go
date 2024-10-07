package trustedaccess

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrustedAccessClient struct {
	Client *resourcemanager.Client
}

func NewTrustedAccessClientWithBaseURI(sdkApi sdkEnv.Api) (*TrustedAccessClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "trustedaccess", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TrustedAccessClient: %+v", err)
	}

	return &TrustedAccessClient{
		Client: client,
	}, nil
}
