package databases

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByServerOperationResponse struct {
	HttpResponse *http.Response
	Model        *DatabaseListResult
}

// ListByServer ...
func (c DatabasesClient) ListByServer(ctx context.Context, id ServerId) (result ListByServerOperationResponse, err error) {
	req, err := c.preparerForListByServer(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "ListByServer", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "ListByServer", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListByServer(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "ListByServer", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListByServer prepares the ListByServer request.
func (c DatabasesClient) preparerForListByServer(ctx context.Context, id ServerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/databases", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListByServer handles the response to the ListByServer request. The method always
// closes the http.Response Body.
func (c DatabasesClient) responderForListByServer(resp *http.Response) (result ListByServerOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
