package applicationinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbookTemplatesUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *WorkbookTemplate
}

// WorkbookTemplatesUpdate ...
func (c ApplicationInsightsClient) WorkbookTemplatesUpdate(ctx context.Context, id WorkbookTemplateId, input WorkbookTemplateUpdateParameters) (result WorkbookTemplatesUpdateOperationResponse, err error) {
	req, err := c.preparerForWorkbookTemplatesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbookTemplatesUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbookTemplatesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkbookTemplatesUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbookTemplatesUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkbookTemplatesUpdate prepares the WorkbookTemplatesUpdate request.
func (c ApplicationInsightsClient) preparerForWorkbookTemplatesUpdate(ctx context.Context, id WorkbookTemplateId, input WorkbookTemplateUpdateParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWorkbookTemplatesUpdate handles the response to the WorkbookTemplatesUpdate request. The method always
// closes the http.Response Body.
func (c ApplicationInsightsClient) responderForWorkbookTemplatesUpdate(resp *http.Response) (result WorkbookTemplatesUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
