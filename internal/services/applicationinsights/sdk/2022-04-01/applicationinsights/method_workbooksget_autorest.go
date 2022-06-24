package applicationinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbooksGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Workbook
}

type WorkbooksGetOperationOptions struct {
	CanFetchContent *bool
}

func DefaultWorkbooksGetOperationOptions() WorkbooksGetOperationOptions {
	return WorkbooksGetOperationOptions{}
}

func (o WorkbooksGetOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o WorkbooksGetOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.CanFetchContent != nil {
		out["canFetchContent"] = *o.CanFetchContent
	}

	return out
}

// WorkbooksGet ...
func (c ApplicationInsightsClient) WorkbooksGet(ctx context.Context, id WorkbookId, options WorkbooksGetOperationOptions) (result WorkbooksGetOperationResponse, err error) {
	req, err := c.preparerForWorkbooksGet(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkbooksGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkbooksGet prepares the WorkbooksGet request.
func (c ApplicationInsightsClient) preparerForWorkbooksGet(ctx context.Context, id WorkbookId, options WorkbooksGetOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWorkbooksGet handles the response to the WorkbooksGet request. The method always
// closes the http.Response Body.
func (c ApplicationInsightsClient) responderForWorkbooksGet(resp *http.Response) (result WorkbooksGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
