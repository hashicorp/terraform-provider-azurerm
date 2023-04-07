package views

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteByScopeOperationResponse struct {
	HttpResponse *http.Response
}

// DeleteByScope ...
func (c ViewsClient) DeleteByScope(ctx context.Context, id ScopedViewId) (result DeleteByScopeOperationResponse, err error) {
	req, err := c.preparerForDeleteByScope(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "views.ViewsClient", "DeleteByScope", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "views.ViewsClient", "DeleteByScope", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDeleteByScope(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "views.ViewsClient", "DeleteByScope", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDeleteByScope prepares the DeleteByScope request.
func (c ViewsClient) preparerForDeleteByScope(ctx context.Context, id ScopedViewId) (*http.Request, error) {
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

// responderForDeleteByScope handles the response to the DeleteByScope request. The method always
// closes the http.Response Body.
func (c ViewsClient) responderForDeleteByScope(resp *http.Response) (result DeleteByScopeOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
