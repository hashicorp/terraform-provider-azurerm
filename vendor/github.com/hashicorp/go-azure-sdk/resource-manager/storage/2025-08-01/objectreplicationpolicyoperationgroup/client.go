package objectreplicationpolicyoperationgroup

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ObjectReplicationPolicyOperationGroupClient struct {
	Client *resourcemanager.Client
}

func NewObjectReplicationPolicyOperationGroupClientWithBaseURI(sdkApi sdkEnv.Api) (*ObjectReplicationPolicyOperationGroupClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "objectreplicationpolicyoperationgroup", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ObjectReplicationPolicyOperationGroupClient: %+v", err)
	}

	return &ObjectReplicationPolicyOperationGroupClient{
		Client: client,
	}, nil
}
