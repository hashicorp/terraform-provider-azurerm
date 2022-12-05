package streamingpoliciesandstreaminglocators

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingLocatorsListContentKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *ListContentKeysResponse
}

// StreamingLocatorsListContentKeys ...
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingLocatorsListContentKeys(ctx context.Context, id StreamingLocatorId) (result StreamingLocatorsListContentKeysOperationResponse, err error) {
	req, err := c.preparerForStreamingLocatorsListContentKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsListContentKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsListContentKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForStreamingLocatorsListContentKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsListContentKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForStreamingLocatorsListContentKeys prepares the StreamingLocatorsListContentKeys request.
func (c StreamingPoliciesAndStreamingLocatorsClient) preparerForStreamingLocatorsListContentKeys(ctx context.Context, id StreamingLocatorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listContentKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForStreamingLocatorsListContentKeys handles the response to the StreamingLocatorsListContentKeys request. The method always
// closes the http.Response Body.
func (c StreamingPoliciesAndStreamingLocatorsClient) responderForStreamingLocatorsListContentKeys(resp *http.Response) (result StreamingLocatorsListContentKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
