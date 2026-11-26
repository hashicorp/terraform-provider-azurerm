package storagetaskassignments

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageTaskAssignmentsClient struct {
	Client *resourcemanager.Client
}

func NewStorageTaskAssignmentsClientWithBaseURI(sdkApi sdkEnv.Api) (*StorageTaskAssignmentsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "storagetaskassignments", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating StorageTaskAssignmentsClient: %+v", err)
	}

	return &StorageTaskAssignmentsClient{
		Client: client,
	}, nil
}
