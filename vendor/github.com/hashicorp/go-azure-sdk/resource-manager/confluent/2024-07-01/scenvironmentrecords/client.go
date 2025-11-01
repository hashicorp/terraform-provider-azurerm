package scenvironmentrecords

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SCEnvironmentRecordsClient struct {
	Client *resourcemanager.Client
}

func NewSCEnvironmentRecordsClientWithBaseURI(sdkApi sdkEnv.Api) (*SCEnvironmentRecordsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "scenvironmentrecords", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating SCEnvironmentRecordsClient: %+v", err)
	}

	return &SCEnvironmentRecordsClient{
		Client: client,
	}, nil
}
