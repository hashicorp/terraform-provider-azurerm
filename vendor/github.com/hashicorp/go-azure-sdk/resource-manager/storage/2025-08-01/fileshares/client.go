package fileshares

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileSharesClient struct {
	Client *resourcemanager.Client
}

func NewFileSharesClientWithBaseURI(sdkApi sdkEnv.Api) (*FileSharesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "fileshares", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FileSharesClient: %+v", err)
	}

	return &FileSharesClient{
		Client: client,
	}, nil
}
