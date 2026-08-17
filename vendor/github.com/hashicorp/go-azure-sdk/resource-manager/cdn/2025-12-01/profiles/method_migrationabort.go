package profiles

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

type MigrationAbortOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// MigrationAbort ...
func (c ProfilesClient) MigrationAbort(ctx context.Context, id ProfileId) (result MigrationAbortOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/migrationAbort", id.ID()),
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

// MigrationAbortThenPoll performs MigrationAbort then polls until it's completed
func (c ProfilesClient) MigrationAbortThenPoll(ctx context.Context, id ProfileId) error {
	return c.MigrationAbortCallbackThenPoll(ctx, id, nil)
}

// MigrationAbortCallbackThenPoll performs MigrationAbort, runs the optional callback function, then polls until it's completed
func (c ProfilesClient) MigrationAbortCallbackThenPoll(ctx context.Context, id ProfileId, callback func() error) error {
	result, err := c.MigrationAbort(ctx, id)
	if err != nil {
		return fmt.Errorf("performing MigrationAbort: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after MigrationAbort: %+v", err)
	}

	return nil
}
