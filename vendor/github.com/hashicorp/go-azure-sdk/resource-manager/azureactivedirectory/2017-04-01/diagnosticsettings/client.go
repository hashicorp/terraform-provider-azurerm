package diagnosticsettings

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticSettingsClient struct {
	Client *resourcemanager.Client
}

func NewDiagnosticSettingsClientWithBaseURI(sdkApi sdkEnv.Api) (*DiagnosticSettingsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "diagnosticsettings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DiagnosticSettingsClient: %+v", err)
	}

	return &DiagnosticSettingsClient{
		Client: client,
	}, nil
}
