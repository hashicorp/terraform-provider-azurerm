package exports

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetExecutionHistoryOperationResponse struct {
	HttpResponse *http.Response
	Model        *ExportExecutionListResult
}

// GetExecutionHistory ...
func (c ExportsClient) GetExecutionHistory(ctx context.Context, id ScopedExportId) (result GetExecutionHistoryOperationResponse, err error) {
	req, err := c.preparerForGetExecutionHistory(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "exports.ExportsClient", "GetExecutionHistory", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "exports.ExportsClient", "GetExecutionHistory", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetExecutionHistory(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "exports.ExportsClient", "GetExecutionHistory", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetExecutionHistory prepares the GetExecutionHistory request.
func (c ExportsClient) preparerForGetExecutionHistory(ctx context.Context, id ScopedExportId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/runHistory", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetExecutionHistory handles the response to the GetExecutionHistory request. The method always
// closes the http.Response Body.
func (c ExportsClient) responderForGetExecutionHistory(resp *http.Response) (result GetExecutionHistoryOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
