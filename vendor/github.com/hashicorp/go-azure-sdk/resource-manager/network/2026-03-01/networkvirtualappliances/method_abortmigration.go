package networkvirtualappliances

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

type AbortMigrationOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// AbortMigration ...
func (c NetworkVirtualAppliancesClient) AbortMigration(ctx context.Context, id NetworkVirtualApplianceId) (result AbortMigrationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/abortMigration", id.ID()),
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

// AbortMigrationThenPoll performs AbortMigration then polls until it's completed
func (c NetworkVirtualAppliancesClient) AbortMigrationThenPoll(ctx context.Context, id NetworkVirtualApplianceId) error {
	return c.AbortMigrationCallbackThenPoll(ctx, id, nil)
}

// AbortMigrationCallbackThenPoll performs AbortMigration, runs the optional callback function, then polls until it's completed
func (c NetworkVirtualAppliancesClient) AbortMigrationCallbackThenPoll(ctx context.Context, id NetworkVirtualApplianceId, callback func() error) error {
	result, err := c.AbortMigration(ctx, id)
	if err != nil {
		return fmt.Errorf("performing AbortMigration: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after AbortMigration: %+v", err)
	}

	return nil
}
