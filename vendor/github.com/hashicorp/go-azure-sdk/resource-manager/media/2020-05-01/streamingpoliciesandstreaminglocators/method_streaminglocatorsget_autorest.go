package streamingpoliciesandstreaminglocators

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingLocatorsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *StreamingLocator
}

// StreamingLocatorsGet ...
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingLocatorsGet(ctx context.Context, id StreamingLocatorId) (result StreamingLocatorsGetOperationResponse, err error) {
	req, err := c.preparerForStreamingLocatorsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForStreamingLocatorsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForStreamingLocatorsGet prepares the StreamingLocatorsGet request.
func (c StreamingPoliciesAndStreamingLocatorsClient) preparerForStreamingLocatorsGet(ctx context.Context, id StreamingLocatorId) (*http.Request, error) {
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

// responderForStreamingLocatorsGet handles the response to the StreamingLocatorsGet request. The method always
// closes the http.Response Body.
func (c StreamingPoliciesAndStreamingLocatorsClient) responderForStreamingLocatorsGet(resp *http.Response) (result StreamingLocatorsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
