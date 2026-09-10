package fullrestore

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/dataplane"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RestoreOperation
}

// Operation ...
func (c FullRestoreClient) Operation(ctx context.Context, input RestoreOperationParameters) (result OperationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodPut,
		Path:       "/restore",
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

	result.Poller, err = dataplane.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// OperationThenPoll performs Operation then polls until it's completed
func (c FullRestoreClient) OperationThenPoll(ctx context.Context, input RestoreOperationParameters) error {
	return c.OperationCallbackThenPoll(ctx, input, nil)
}

// OperationCallbackThenPoll performs Operation, runs the optional callback function, then polls until it's completed
func (c FullRestoreClient) OperationCallbackThenPoll(ctx context.Context, input RestoreOperationParameters, callback func() error) error {
	result, err := c.Operation(ctx, input)
	if err != nil {
		return fmt.Errorf("performing Operation: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Operation: %+v", err)
	}

	return nil
}
