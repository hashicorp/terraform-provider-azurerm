package vminstancehybrididentitymetadatas

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMInstanceHybridIdentityMetadatasClient struct {
	Client *resourcemanager.Client
}

func NewVMInstanceHybridIdentityMetadatasClientWithBaseURI(sdkApi sdkEnv.Api) (*VMInstanceHybridIdentityMetadatasClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "vminstancehybrididentitymetadatas", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VMInstanceHybridIdentityMetadatasClient: %+v", err)
	}

	return &VMInstanceHybridIdentityMetadatasClient{
		Client: client,
	}, nil
}
