package recordsets

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecordSetsClient struct {
	Client *resourcemanager.Client
}

func NewRecordSetsClientWithBaseURI(sdkApi sdkEnv.Api) (*RecordSetsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "recordsets", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RecordSetsClient: %+v", err)
	}

	return &RecordSetsClient{
		Client: client,
	}, nil
}
