package applicationinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbooksCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Workbook
}

type WorkbooksCreateOrUpdateOperationOptions struct {
	SourceId *string
}

func DefaultWorkbooksCreateOrUpdateOperationOptions() WorkbooksCreateOrUpdateOperationOptions {
	return WorkbooksCreateOrUpdateOperationOptions{}
}

func (o WorkbooksCreateOrUpdateOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o WorkbooksCreateOrUpdateOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.SourceId != nil {
		out["sourceId"] = *o.SourceId
	}

	return out
}

// WorkbooksCreateOrUpdate ...
func (c ApplicationInsightsClient) WorkbooksCreateOrUpdate(ctx context.Context, id WorkbookId, input Workbook, options WorkbooksCreateOrUpdateOperationOptions) (result WorkbooksCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForWorkbooksCreateOrUpdate(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkbooksCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkbooksCreateOrUpdate prepares the WorkbooksCreateOrUpdate request.
func (c ApplicationInsightsClient) preparerForWorkbooksCreateOrUpdate(ctx context.Context, id WorkbookId, input Workbook, options WorkbooksCreateOrUpdateOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWorkbooksCreateOrUpdate handles the response to the WorkbooksCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c ApplicationInsightsClient) responderForWorkbooksCreateOrUpdate(resp *http.Response) (result WorkbooksCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
