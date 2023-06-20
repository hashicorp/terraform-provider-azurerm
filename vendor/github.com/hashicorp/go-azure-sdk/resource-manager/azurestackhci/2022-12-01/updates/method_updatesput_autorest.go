package updates

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdatesPutOperationResponse struct {
	HttpResponse *http.Response
	Model        *Update
}

// UpdatesPut ...
func (c UpdatesClient) UpdatesPut(ctx context.Context, id UpdateId, input Update) (result UpdatesPutOperationResponse, err error) {
	req, err := c.preparerForUpdatesPut(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updates.UpdatesClient", "UpdatesPut", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "updates.UpdatesClient", "UpdatesPut", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUpdatesPut(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updates.UpdatesClient", "UpdatesPut", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUpdatesPut prepares the UpdatesPut request.
func (c UpdatesClient) preparerForUpdatesPut(ctx context.Context, id UpdateId, input Update) (*http.Request, error) {
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

// responderForUpdatesPut handles the response to the UpdatesPut request. The method always
// closes the http.Response Body.
func (c UpdatesClient) responderForUpdatesPut(resp *http.Response) (result UpdatesPutOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
