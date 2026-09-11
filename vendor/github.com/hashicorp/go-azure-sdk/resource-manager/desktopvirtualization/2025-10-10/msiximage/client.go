package msiximage

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MsixImageClient struct {
	Client *resourcemanager.Client
}

func NewMsixImageClientWithBaseURI(sdkApi sdkEnv.Api) (*MsixImageClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "msiximage", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating MsixImageClient: %+v", err)
	}

	return &MsixImageClient{
		Client: client,
	}, nil
}
