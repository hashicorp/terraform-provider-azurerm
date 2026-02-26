package namedvalue

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceNamedValueRefreshSecretOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *NamedValueContract
}

// WorkspaceNamedValueRefreshSecret ...
func (c NamedValueClient) WorkspaceNamedValueRefreshSecret(ctx context.Context, id WorkspaceNamedValueId) (result WorkspaceNamedValueRefreshSecretOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/refreshSecret", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// WorkspaceNamedValueRefreshSecretThenPoll performs WorkspaceNamedValueRefreshSecret then polls until it's completed
func (c NamedValueClient) WorkspaceNamedValueRefreshSecretThenPoll(ctx context.Context, id WorkspaceNamedValueId) error {
	result, err := c.WorkspaceNamedValueRefreshSecret(ctx, id)
	if err != nil {
		return fmt.Errorf("performing WorkspaceNamedValueRefreshSecret: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after WorkspaceNamedValueRefreshSecret: %+v", err)
	}

	return nil
}
