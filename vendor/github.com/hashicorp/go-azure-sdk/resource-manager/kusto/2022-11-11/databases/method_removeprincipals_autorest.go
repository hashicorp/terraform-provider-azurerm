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

type RemovePrincipalsOperationResponse struct {
	HttpResponse *http.Response
	Model        *DatabasePrincipalListResult
}

// RemovePrincipals ...
func (c DatabasesClient) RemovePrincipals(ctx context.Context, id DatabaseId, input DatabasePrincipalListRequest) (result RemovePrincipalsOperationResponse, err error) {
	req, err := c.preparerForRemovePrincipals(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "RemovePrincipals", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "RemovePrincipals", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRemovePrincipals(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "RemovePrincipals", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRemovePrincipals prepares the RemovePrincipals request.
func (c DatabasesClient) preparerForRemovePrincipals(ctx context.Context, id DatabaseId, input DatabasePrincipalListRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/removePrincipals", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRemovePrincipals handles the response to the RemovePrincipals request. The method always
// closes the http.Response Body.
func (c DatabasesClient) responderForRemovePrincipals(resp *http.Response) (result RemovePrincipalsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
