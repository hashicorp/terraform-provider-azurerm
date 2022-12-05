package streamingpoliciesandstreaminglocators

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingPoliciesDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// StreamingPoliciesDelete ...
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingPoliciesDelete(ctx context.Context, id StreamingPolicyId) (result StreamingPoliciesDeleteOperationResponse, err error) {
	req, err := c.preparerForStreamingPoliciesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingPoliciesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingPoliciesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForStreamingPoliciesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingPoliciesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForStreamingPoliciesDelete prepares the StreamingPoliciesDelete request.
func (c StreamingPoliciesAndStreamingLocatorsClient) preparerForStreamingPoliciesDelete(ctx context.Context, id StreamingPolicyId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForStreamingPoliciesDelete handles the response to the StreamingPoliciesDelete request. The method always
// closes the http.Response Body.
func (c StreamingPoliciesAndStreamingLocatorsClient) responderForStreamingPoliciesDelete(resp *http.Response) (result StreamingPoliciesDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
