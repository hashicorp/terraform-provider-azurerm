package accounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MediaservicesListEdgePoliciesOperationResponse struct {
	HttpResponse *http.Response
	Model        *EdgePolicies
}

// MediaservicesListEdgePolicies ...
func (c AccountsClient) MediaservicesListEdgePolicies(ctx context.Context, id MediaServiceId, input ListEdgePoliciesInput) (result MediaservicesListEdgePoliciesOperationResponse, err error) {
	req, err := c.preparerForMediaservicesListEdgePolicies(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesListEdgePolicies", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesListEdgePolicies", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMediaservicesListEdgePolicies(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesListEdgePolicies", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMediaservicesListEdgePolicies prepares the MediaservicesListEdgePolicies request.
func (c AccountsClient) preparerForMediaservicesListEdgePolicies(ctx context.Context, id MediaServiceId, input ListEdgePoliciesInput) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listEdgePolicies", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMediaservicesListEdgePolicies handles the response to the MediaservicesListEdgePolicies request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForMediaservicesListEdgePolicies(resp *http.Response) (result MediaservicesListEdgePoliciesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
