package servicelinker

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkerGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *LinkerResource
}

// LinkerGet ...
func (c ServiceLinkerClient) LinkerGet(ctx context.Context, id ScopedLinkerId) (result LinkerGetOperationResponse, err error) {
	req, err := c.preparerForLinkerGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicelinker.ServiceLinkerClient", "LinkerGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicelinker.ServiceLinkerClient", "LinkerGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForLinkerGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicelinker.ServiceLinkerClient", "LinkerGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForLinkerGet prepares the LinkerGet request.
func (c ServiceLinkerClient) preparerForLinkerGet(ctx context.Context, id ScopedLinkerId) (*http.Request, error) {
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

// responderForLinkerGet handles the response to the LinkerGet request. The method always
// closes the http.Response Body.
func (c ServiceLinkerClient) responderForLinkerGet(resp *http.Response) (result LinkerGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
