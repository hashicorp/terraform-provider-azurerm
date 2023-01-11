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

type StreamingLocatorsListPathsOperationResponse struct {
	HttpResponse *http.Response
	Model        *ListPathsResponse
}

// StreamingLocatorsListPaths ...
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingLocatorsListPaths(ctx context.Context, id StreamingLocatorId) (result StreamingLocatorsListPathsOperationResponse, err error) {
	req, err := c.preparerForStreamingLocatorsListPaths(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsListPaths", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsListPaths", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForStreamingLocatorsListPaths(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsListPaths", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForStreamingLocatorsListPaths prepares the StreamingLocatorsListPaths request.
func (c StreamingPoliciesAndStreamingLocatorsClient) preparerForStreamingLocatorsListPaths(ctx context.Context, id StreamingLocatorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listPaths", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForStreamingLocatorsListPaths handles the response to the StreamingLocatorsListPaths request. The method always
// closes the http.Response Body.
func (c StreamingPoliciesAndStreamingLocatorsClient) responderForStreamingLocatorsListPaths(resp *http.Response) (result StreamingLocatorsListPathsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
