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

type ScriptsListByDatabaseOperationResponse struct {
	HttpResponse *http.Response
	Model        *ScriptListResult
}

// ScriptsListByDatabase ...
func (c ScriptsClient) ScriptsListByDatabase(ctx context.Context, id DatabaseId) (result ScriptsListByDatabaseOperationResponse, err error) {
	req, err := c.preparerForScriptsListByDatabase(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scripts.ScriptsClient", "ScriptsListByDatabase", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "scripts.ScriptsClient", "ScriptsListByDatabase", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForScriptsListByDatabase(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scripts.ScriptsClient", "ScriptsListByDatabase", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForScriptsListByDatabase prepares the ScriptsListByDatabase request.
func (c ScriptsClient) preparerForScriptsListByDatabase(ctx context.Context, id DatabaseId) (*http.Request, error) {
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

// responderForScriptsListByDatabase handles the response to the ScriptsListByDatabase request. The method always
// closes the http.Response Body.
func (c ScriptsClient) responderForScriptsListByDatabase(resp *http.Response) (result ScriptsListByDatabaseOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
