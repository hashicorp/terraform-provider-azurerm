package extensions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtensionsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Extension
}

// ExtensionsGet ...
func (c ExtensionsClient) ExtensionsGet(ctx context.Context, id ExtensionId) (result ExtensionsGetOperationResponse, err error) {
	req, err := c.preparerForExtensionsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "extensions.ExtensionsClient", "ExtensionsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "extensions.ExtensionsClient", "ExtensionsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForExtensionsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "extensions.ExtensionsClient", "ExtensionsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForExtensionsGet prepares the ExtensionsGet request.
func (c ExtensionsClient) preparerForExtensionsGet(ctx context.Context, id ExtensionId) (*http.Request, error) {
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

// responderForExtensionsGet handles the response to the ExtensionsGet request. The method always
// closes the http.Response Body.
func (c ExtensionsClient) responderForExtensionsGet(resp *http.Response) (result ExtensionsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
