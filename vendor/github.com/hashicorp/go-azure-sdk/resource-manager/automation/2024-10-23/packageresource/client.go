package packageresource

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PackageResourceClient struct {
	Client *resourcemanager.Client
}

func NewPackageResourceClientWithBaseURI(sdkApi sdkEnv.Api) (*PackageResourceClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "packageresource", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PackageResourceClient: %+v", err)
	}

	return &PackageResourceClient{
		Client: client,
	}, nil
}
