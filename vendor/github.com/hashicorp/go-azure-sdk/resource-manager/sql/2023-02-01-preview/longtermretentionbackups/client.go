package longtermretentionbackups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LongTermRetentionBackupsClient struct {
	Client *resourcemanager.Client
}

func NewLongTermRetentionBackupsClientWithBaseURI(sdkApi sdkEnv.Api) (*LongTermRetentionBackupsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "longtermretentionbackups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LongTermRetentionBackupsClient: %+v", err)
	}

	return &LongTermRetentionBackupsClient{
		Client: client,
	}, nil
}
