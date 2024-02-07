package diagnostic

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticClient struct {
	Client *resourcemanager.Client
}

func NewDiagnosticClientWithBaseURI(sdkApi sdkEnv.Api) (*DiagnosticClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "diagnostic", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DiagnosticClient: %+v", err)
	}

	return &DiagnosticClient{
		Client: client,
	}, nil
}
