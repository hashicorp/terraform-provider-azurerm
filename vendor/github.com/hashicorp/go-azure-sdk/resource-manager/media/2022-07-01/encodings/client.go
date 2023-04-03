package encodings

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncodingsClient struct {
	Client *resourcemanager.Client
}

func NewEncodingsClientWithBaseURI(api environments.Api) (*EncodingsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "encodings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating EncodingsClient: %+v", err)
	}

	return &EncodingsClient{
		Client: client,
	}, nil
}
