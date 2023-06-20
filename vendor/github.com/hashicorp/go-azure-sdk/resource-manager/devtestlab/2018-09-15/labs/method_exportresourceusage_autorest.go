package labs

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

type ExportResourceUsageOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ExportResourceUsage ...
func (c LabsClient) ExportResourceUsage(ctx context.Context, id LabId, input ExportResourceUsageParameters) (result ExportResourceUsageOperationResponse, err error) {
	req, err := c.preparerForExportResourceUsage(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "labs.LabsClient", "ExportResourceUsage", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForExportResourceUsage(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "labs.LabsClient", "ExportResourceUsage", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ExportResourceUsageThenPoll performs ExportResourceUsage then polls until it's completed
func (c LabsClient) ExportResourceUsageThenPoll(ctx context.Context, id LabId, input ExportResourceUsageParameters) error {
	result, err := c.ExportResourceUsage(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ExportResourceUsage: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ExportResourceUsage: %+v", err)
	}

	return nil
}

// preparerForExportResourceUsage prepares the ExportResourceUsage request.
func (c LabsClient) preparerForExportResourceUsage(ctx context.Context, id LabId, input ExportResourceUsageParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/exportResourceUsage", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForExportResourceUsage sends the ExportResourceUsage request. The method will close the
// http.Response Body if it receives an error.
func (c LabsClient) senderForExportResourceUsage(ctx context.Context, req *http.Request) (future ExportResourceUsageOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
