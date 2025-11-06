package cloudexadatainfrastructures

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudExadataInfrastructuresClient struct {
	Client *resourcemanager.Client
}

func NewCloudExadataInfrastructuresClientWithBaseURI(sdkApi sdkEnv.Api) (*CloudExadataInfrastructuresClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "cloudexadatainfrastructures", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CloudExadataInfrastructuresClient: %+v", err)
	}

	return &CloudExadataInfrastructuresClient{
		Client: client,
	}, nil
}
