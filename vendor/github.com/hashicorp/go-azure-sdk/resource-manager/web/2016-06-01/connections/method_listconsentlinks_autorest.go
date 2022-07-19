package connections

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListConsentLinksOperationResponse struct {
	HttpResponse *http.Response
	Model        *ConsentLinkCollection
}

// ListConsentLinks ...
func (c ConnectionsClient) ListConsentLinks(ctx context.Context, id ConnectionId, input ListConsentLinksDefinition) (result ListConsentLinksOperationResponse, err error) {
	req, err := c.preparerForListConsentLinks(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connections.ConnectionsClient", "ListConsentLinks", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "connections.ConnectionsClient", "ListConsentLinks", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListConsentLinks(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connections.ConnectionsClient", "ListConsentLinks", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListConsentLinks prepares the ListConsentLinks request.
func (c ConnectionsClient) preparerForListConsentLinks(ctx context.Context, id ConnectionId, input ListConsentLinksDefinition) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listConsentLinks", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListConsentLinks handles the response to the ListConsentLinks request. The method always
// closes the http.Response Body.
func (c ConnectionsClient) responderForListConsentLinks(resp *http.Response) (result ListConsentLinksOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
