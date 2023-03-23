package workbooktemplatesapis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbookTemplatesDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// WorkbookTemplatesDelete ...
func (c WorkbookTemplatesAPIsClient) WorkbookTemplatesDelete(ctx context.Context, id WorkbookTemplateId) (result WorkbookTemplatesDeleteOperationResponse, err error) {
	req, err := c.preparerForWorkbookTemplatesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workbooktemplatesapis.WorkbookTemplatesAPIsClient", "WorkbookTemplatesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workbooktemplatesapis.WorkbookTemplatesAPIsClient", "WorkbookTemplatesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkbookTemplatesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workbooktemplatesapis.WorkbookTemplatesAPIsClient", "WorkbookTemplatesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkbookTemplatesDelete prepares the WorkbookTemplatesDelete request.
func (c WorkbookTemplatesAPIsClient) preparerForWorkbookTemplatesDelete(ctx context.Context, id WorkbookTemplateId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWorkbookTemplatesDelete handles the response to the WorkbookTemplatesDelete request. The method always
// closes the http.Response Body.
func (c WorkbookTemplatesAPIsClient) responderForWorkbookTemplatesDelete(resp *http.Response) (result WorkbookTemplatesDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
