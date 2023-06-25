package dataconnections

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataConnectionValidationOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DataConnectionValidation ...
func (c DataConnectionsClient) DataConnectionValidation(ctx context.Context, id DatabaseId, input DataConnectionValidation) (result DataConnectionValidationOperationResponse, err error) {
	req, err := c.preparerForDataConnectionValidation(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dataconnections.DataConnectionsClient", "DataConnectionValidation", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDataConnectionValidation(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dataconnections.DataConnectionsClient", "DataConnectionValidation", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DataConnectionValidationThenPoll performs DataConnectionValidation then polls until it's completed
func (c DataConnectionsClient) DataConnectionValidationThenPoll(ctx context.Context, id DatabaseId, input DataConnectionValidation) error {
	result, err := c.DataConnectionValidation(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DataConnectionValidation: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DataConnectionValidation: %+v", err)
	}

	return nil
}

// preparerForDataConnectionValidation prepares the DataConnectionValidation request.
func (c DataConnectionsClient) preparerForDataConnectionValidation(ctx context.Context, id DatabaseId, input DataConnectionValidation) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/dataConnectionValidation", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDataConnectionValidation sends the DataConnectionValidation request. The method will close the
// http.Response Body if it receives an error.
func (c DataConnectionsClient) senderForDataConnectionValidation(ctx context.Context, req *http.Request) (future DataConnectionValidationOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
