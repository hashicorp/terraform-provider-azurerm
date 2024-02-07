package labs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LabsClient struct {
	Client *resourcemanager.Client
}

func NewLabsClientWithBaseURI(sdkApi sdkEnv.Api) (*LabsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "labs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LabsClient: %+v", err)
	}

	return &LabsClient{
		Client: client,
	}, nil
}
