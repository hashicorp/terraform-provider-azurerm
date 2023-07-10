package publishers

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PublishersGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Publisher
}

// PublishersGet ...
func (c PublishersClient) PublishersGet(ctx context.Context, id PublisherId) (result PublishersGetOperationResponse, err error) {
	req, err := c.preparerForPublishersGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "publishers.PublishersClient", "PublishersGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "publishers.PublishersClient", "PublishersGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPublishersGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "publishers.PublishersClient", "PublishersGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPublishersGet prepares the PublishersGet request.
func (c PublishersClient) preparerForPublishersGet(ctx context.Context, id PublisherId) (*http.Request, error) {
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

// responderForPublishersGet handles the response to the PublishersGet request. The method always
// closes the http.Response Body.
func (c PublishersClient) responderForPublishersGet(resp *http.Response) (result PublishersGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
