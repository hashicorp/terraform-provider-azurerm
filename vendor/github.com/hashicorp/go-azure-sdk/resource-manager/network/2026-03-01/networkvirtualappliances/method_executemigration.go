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

type ExecuteMigrationOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// ExecuteMigration ...
func (c NetworkVirtualAppliancesClient) ExecuteMigration(ctx context.Context, id NetworkVirtualApplianceId, input NetworkVirtualApplianceExecuteMigrationRequest) (result ExecuteMigrationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusNoContent,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/executeMigration", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
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

// ExecuteMigrationThenPoll performs ExecuteMigration then polls until it's completed
func (c NetworkVirtualAppliancesClient) ExecuteMigrationThenPoll(ctx context.Context, id NetworkVirtualApplianceId, input NetworkVirtualApplianceExecuteMigrationRequest) error {
	return c.ExecuteMigrationCallbackThenPoll(ctx, id, input, nil)
}

// ExecuteMigrationCallbackThenPoll performs ExecuteMigration, runs the optional callback function, then polls until it's completed
func (c NetworkVirtualAppliancesClient) ExecuteMigrationCallbackThenPoll(ctx context.Context, id NetworkVirtualApplianceId, input NetworkVirtualApplianceExecuteMigrationRequest, callback func() error) error {
	result, err := c.ExecuteMigration(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ExecuteMigration: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ExecuteMigration: %+v", err)
	}

	return nil
}
