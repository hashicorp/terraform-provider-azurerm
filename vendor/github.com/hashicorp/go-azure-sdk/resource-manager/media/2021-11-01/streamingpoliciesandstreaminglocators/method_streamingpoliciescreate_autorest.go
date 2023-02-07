package streamingpoliciesandstreaminglocators

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingPoliciesCreateOperationResponse struct {
	HttpResponse *http.Response
	Model        *StreamingPolicy
}

// StreamingPoliciesCreate ...
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingPoliciesCreate(ctx context.Context, id StreamingPolicyId, input StreamingPolicy) (result StreamingPoliciesCreateOperationResponse, err error) {
	req, err := c.preparerForStreamingPoliciesCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingPoliciesCreate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingPoliciesCreate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForStreamingPoliciesCreate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingPoliciesCreate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForStreamingPoliciesCreate prepares the StreamingPoliciesCreate request.
func (c StreamingPoliciesAndStreamingLocatorsClient) preparerForStreamingPoliciesCreate(ctx context.Context, id StreamingPolicyId, input StreamingPolicy) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForStreamingPoliciesCreate handles the response to the StreamingPoliciesCreate request. The method always
// closes the http.Response Body.
func (c StreamingPoliciesAndStreamingLocatorsClient) responderForStreamingPoliciesCreate(resp *http.Response) (result StreamingPoliciesCreateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
