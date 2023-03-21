package scripts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByDatabaseOperationResponse struct {
	HttpResponse *http.Response
	Model        *ScriptListResult
}

// ListByDatabase ...
func (c ScriptsClient) ListByDatabase(ctx context.Context, id DatabaseId) (result ListByDatabaseOperationResponse, err error) {
	req, err := c.preparerForListByDatabase(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scripts.ScriptsClient", "ListByDatabase", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "scripts.ScriptsClient", "ListByDatabase", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListByDatabase(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scripts.ScriptsClient", "ListByDatabase", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListByDatabase prepares the ListByDatabase request.
func (c ScriptsClient) preparerForListByDatabase(ctx context.Context, id DatabaseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/scripts", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListByDatabase handles the response to the ListByDatabase request. The method always
// closes the http.Response Body.
func (c ScriptsClient) responderForListByDatabase(resp *http.Response) (result ListByDatabaseOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
