package applicationinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbookTemplatesCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *WorkbookTemplate
}

// WorkbookTemplatesCreateOrUpdate ...
func (c ApplicationInsightsClient) WorkbookTemplatesCreateOrUpdate(ctx context.Context, id WorkbookTemplateId, input WorkbookTemplate) (result WorkbookTemplatesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForWorkbookTemplatesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbookTemplatesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbookTemplatesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkbookTemplatesCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbookTemplatesCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkbookTemplatesCreateOrUpdate prepares the WorkbookTemplatesCreateOrUpdate request.
func (c ApplicationInsightsClient) preparerForWorkbookTemplatesCreateOrUpdate(ctx context.Context, id WorkbookTemplateId, input WorkbookTemplate) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWorkbookTemplatesCreateOrUpdate handles the response to the WorkbookTemplatesCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c ApplicationInsightsClient) responderForWorkbookTemplatesCreateOrUpdate(resp *http.Response) (result WorkbookTemplatesCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
