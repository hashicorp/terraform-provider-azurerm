package scriptexecutionhistory

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScriptExecutionHistoryClient struct {
	Client *resourcemanager.Client
}

func NewScriptExecutionHistoryClientWithBaseURI(sdkApi sdkEnv.Api) (*ScriptExecutionHistoryClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "scriptexecutionhistory", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ScriptExecutionHistoryClient: %+v", err)
	}

	return &ScriptExecutionHistoryClient{
		Client: client,
	}, nil
}
