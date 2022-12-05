package streamingpoliciesandstreaminglocators

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingLocatorsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// StreamingLocatorsDelete ...
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingLocatorsDelete(ctx context.Context, id StreamingLocatorId) (result StreamingLocatorsDeleteOperationResponse, err error) {
	req, err := c.preparerForStreamingLocatorsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForStreamingLocatorsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForStreamingLocatorsDelete prepares the StreamingLocatorsDelete request.
func (c StreamingPoliciesAndStreamingLocatorsClient) preparerForStreamingLocatorsDelete(ctx context.Context, id StreamingLocatorId) (*http.Request, error) {
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

// responderForStreamingLocatorsDelete handles the response to the StreamingLocatorsDelete request. The method always
// closes the http.Response Body.
func (c StreamingPoliciesAndStreamingLocatorsClient) responderForStreamingLocatorsDelete(resp *http.Response) (result StreamingLocatorsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
