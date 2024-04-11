package transformations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TransformationsClient struct {
	Client *resourcemanager.Client
}

func NewTransformationsClientWithBaseURI(sdkApi sdkEnv.Api) (*TransformationsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "transformations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating TransformationsClient: %+v", err)
	}

	return &TransformationsClient{
		Client: client,
	}, nil
}
