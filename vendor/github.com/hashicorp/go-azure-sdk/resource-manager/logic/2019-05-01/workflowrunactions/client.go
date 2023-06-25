package workflowrunactions

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkflowRunActionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewWorkflowRunActionsClientWithBaseURI(endpoint string) WorkflowRunActionsClient {
	return WorkflowRunActionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}
