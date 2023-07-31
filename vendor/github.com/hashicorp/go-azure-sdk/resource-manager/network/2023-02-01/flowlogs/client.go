package flowlogs

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FlowLogsClient struct {
	Client *resourcemanager.Client
}

func NewFlowLogsClientWithBaseURI(api environments.Api) (*FlowLogsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "flowlogs", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FlowLogsClient: %+v", err)
	}

	return &FlowLogsClient{
		Client: client,
	}, nil
}
