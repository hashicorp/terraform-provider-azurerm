package dataexport

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataExportClient struct {
	Client *resourcemanager.Client
}

func NewDataExportClientWithBaseURI(sdkApi sdkEnv.Api) (*DataExportClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "dataexport", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DataExportClient: %+v", err)
	}

	return &DataExportClient{
		Client: client,
	}, nil
}
