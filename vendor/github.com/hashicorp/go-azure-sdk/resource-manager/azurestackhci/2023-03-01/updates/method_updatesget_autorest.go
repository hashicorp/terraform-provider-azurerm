package updates

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdatesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Update
}

// UpdatesGet ...
func (c UpdatesClient) UpdatesGet(ctx context.Context, id UpdateId) (result UpdatesGetOperationResponse, err error) {
	req, err := c.preparerForUpdatesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updates.UpdatesClient", "UpdatesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "updates.UpdatesClient", "UpdatesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUpdatesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updates.UpdatesClient", "UpdatesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUpdatesGet prepares the UpdatesGet request.
func (c UpdatesClient) preparerForUpdatesGet(ctx context.Context, id UpdateId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUpdatesGet handles the response to the UpdatesGet request. The method always
// closes the http.Response Body.
func (c UpdatesClient) responderForUpdatesGet(resp *http.Response) (result UpdatesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
