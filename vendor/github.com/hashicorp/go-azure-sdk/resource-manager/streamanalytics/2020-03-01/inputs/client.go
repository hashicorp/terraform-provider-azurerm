package inputs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InputsClient struct {
	Client *resourcemanager.Client
}

func NewInputsClientWithBaseURI(sdkApi sdkEnv.Api) (*InputsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "inputs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating InputsClient: %+v", err)
	}

	return &InputsClient{
		Client: client,
	}, nil
}
