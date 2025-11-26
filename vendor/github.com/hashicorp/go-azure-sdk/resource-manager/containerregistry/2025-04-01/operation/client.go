package operation

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationClient struct {
	Client *resourcemanager.Client
}

func NewOperationClientWithBaseURI(sdkApi sdkEnv.Api) (*OperationClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "operation", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating OperationClient: %+v", err)
	}

	return &OperationClient{
		Client: client,
	}, nil
}
