package contactprofile

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContactProfileClient struct {
	Client *resourcemanager.Client
}

func NewContactProfileClientWithBaseURI(api environments.Api) (*ContactProfileClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "contactprofile", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ContactProfileClient: %+v", err)
	}

	return &ContactProfileClient{
		Client: client,
	}, nil
}
