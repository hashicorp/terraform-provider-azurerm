package blobservice

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobServiceClient struct {
	Client *resourcemanager.Client
}

func NewBlobServiceClientWithBaseURI(sdkApi sdkEnv.Api) (*BlobServiceClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "blobservice", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BlobServiceClient: %+v", err)
	}

	return &BlobServiceClient{
		Client: client,
	}, nil
}
