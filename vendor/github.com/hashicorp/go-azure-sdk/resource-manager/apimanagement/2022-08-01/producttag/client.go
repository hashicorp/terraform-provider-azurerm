package producttag

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProductTagClient struct {
	Client *resourcemanager.Client
}

func NewProductTagClientWithBaseURI(sdkApi sdkEnv.Api) (*ProductTagClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "producttag", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ProductTagClient: %+v", err)
	}

	return &ProductTagClient{
		Client: client,
	}, nil
}
