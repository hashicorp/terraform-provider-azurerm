package workbooksapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbooksUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Workbook
}

type WorkbooksUpdateOperationOptions struct {
	SourceId *string
}

func DefaultWorkbooksUpdateOperationOptions() WorkbooksUpdateOperationOptions {
	return WorkbooksUpdateOperationOptions{}
}

func (o WorkbooksUpdateOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o WorkbooksUpdateOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.SourceId != nil {
		out["sourceId"] = *o.SourceId
	}

	return out
}

// WorkbooksUpdate ...
func (c WorkbooksAPIsClient) WorkbooksUpdate(ctx context.Context, id WorkbookId, input WorkbookUpdateParameters, options WorkbooksUpdateOperationOptions) (result WorkbooksUpdateOperationResponse, err error) {
	req, err := c.preparerForWorkbooksUpdate(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workbooksapis.WorkbooksAPIsClient", "WorkbooksUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workbooksapis.WorkbooksAPIsClient", "WorkbooksUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkbooksUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workbooksapis.WorkbooksAPIsClient", "WorkbooksUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkbooksUpdate prepares the WorkbooksUpdate request.
func (c WorkbooksAPIsClient) preparerForWorkbooksUpdate(ctx context.Context, id WorkbookId, input WorkbookUpdateParameters, options WorkbooksUpdateOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWorkbooksUpdate handles the response to the WorkbooksUpdate request. The method always
// closes the http.Response Body.
func (c WorkbooksAPIsClient) responderForWorkbooksUpdate(resp *http.Response) (result WorkbooksUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
