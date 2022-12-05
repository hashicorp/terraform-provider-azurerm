package streamingpoliciesandstreaminglocators

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingPoliciesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *StreamingPolicy
}

// StreamingPoliciesGet ...
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingPoliciesGet(ctx context.Context, id StreamingPolicyId) (result StreamingPoliciesGetOperationResponse, err error) {
	req, err := c.preparerForStreamingPoliciesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingPoliciesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingPoliciesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForStreamingPoliciesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingPoliciesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForStreamingPoliciesGet prepares the StreamingPoliciesGet request.
func (c StreamingPoliciesAndStreamingLocatorsClient) preparerForStreamingPoliciesGet(ctx context.Context, id StreamingPolicyId) (*http.Request, error) {
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

// responderForStreamingPoliciesGet handles the response to the StreamingPoliciesGet request. The method always
// closes the http.Response Body.
func (c StreamingPoliciesAndStreamingLocatorsClient) responderForStreamingPoliciesGet(resp *http.Response) (result StreamingPoliciesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
