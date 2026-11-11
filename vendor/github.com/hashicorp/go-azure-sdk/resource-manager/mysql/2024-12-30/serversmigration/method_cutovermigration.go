package serversmigration

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

type CutoverMigrationOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *Server
}

// CutoverMigration ...
func (c ServersMigrationClient) CutoverMigration(ctx context.Context, id FlexibleServerId) (result CutoverMigrationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/cutoverMigration", id.ID()),
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

// CutoverMigrationThenPoll performs CutoverMigration then polls until it's completed
func (c ServersMigrationClient) CutoverMigrationThenPoll(ctx context.Context, id FlexibleServerId) error {
	return c.CutoverMigrationCallbackThenPoll(ctx, id, nil)
}

// CutoverMigrationCallbackThenPoll performs CutoverMigration, runs the optional callback function, then polls until it's completed
func (c ServersMigrationClient) CutoverMigrationCallbackThenPoll(ctx context.Context, id FlexibleServerId, callback func() error) error {
	result, err := c.CutoverMigration(ctx, id)
	if err != nil {
		return fmt.Errorf("performing CutoverMigration: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after CutoverMigration: %+v", err)
	}

	return nil
}
