package basebackuppolicyresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BaseBackupPolicyResourcesClient struct {
	Client *resourcemanager.Client
}

func NewBaseBackupPolicyResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*BaseBackupPolicyResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "basebackuppolicyresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BaseBackupPolicyResourcesClient: %+v", err)
	}

	return &BaseBackupPolicyResourcesClient{
		Client: client,
	}, nil
}
