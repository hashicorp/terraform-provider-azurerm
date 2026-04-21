package backuppolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewBackupPoliciesClientWithBaseURI(sdkApi sdkEnv.Api) (*BackupPoliciesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "backuppolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BackupPoliciesClient: %+v", err)
	}

	return &BackupPoliciesClient{
		Client: client,
	}, nil
}
