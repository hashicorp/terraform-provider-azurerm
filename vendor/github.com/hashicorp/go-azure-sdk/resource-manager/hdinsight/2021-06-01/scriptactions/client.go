package scriptactions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScriptActionsClient struct {
	Client *resourcemanager.Client
}

func NewScriptActionsClientWithBaseURI(sdkApi sdkEnv.Api) (*ScriptActionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "scriptactions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ScriptActionsClient: %+v", err)
	}

	return &ScriptActionsClient{
		Client: client,
	}, nil
}
