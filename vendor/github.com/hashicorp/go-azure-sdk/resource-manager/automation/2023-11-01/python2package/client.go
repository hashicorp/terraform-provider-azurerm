package python2package

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Python2PackageClient struct {
	Client *resourcemanager.Client
}

func NewPython2PackageClientWithBaseURI(sdkApi sdkEnv.Api) (*Python2PackageClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "python2package", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating Python2PackageClient: %+v", err)
	}

	return &Python2PackageClient{
		Client: client,
	}, nil
}
