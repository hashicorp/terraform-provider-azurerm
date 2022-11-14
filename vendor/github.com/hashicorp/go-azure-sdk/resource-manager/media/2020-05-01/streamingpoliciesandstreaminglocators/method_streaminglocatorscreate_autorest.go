package streamingpoliciesandstreaminglocators

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingLocatorsCreateOperationResponse struct {
	HttpResponse *http.Response
	Model        *StreamingLocator
}

// StreamingLocatorsCreate ...
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingLocatorsCreate(ctx context.Context, id StreamingLocatorId, input StreamingLocator) (result StreamingLocatorsCreateOperationResponse, err error) {
	req, err := c.preparerForStreamingLocatorsCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsCreate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsCreate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForStreamingLocatorsCreate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsCreate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForStreamingLocatorsCreate prepares the StreamingLocatorsCreate request.
func (c StreamingPoliciesAndStreamingLocatorsClient) preparerForStreamingLocatorsCreate(ctx context.Context, id StreamingLocatorId, input StreamingLocator) (*http.Request, error) {
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

// responderForStreamingLocatorsCreate handles the response to the StreamingLocatorsCreate request. The method always
// closes the http.Response Body.
func (c StreamingPoliciesAndStreamingLocatorsClient) responderForStreamingLocatorsCreate(resp *http.Response) (result StreamingLocatorsCreateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
