package topicsauthorizationrule

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicsListKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *AccessKeys
}

// TopicsListKeys ...
func (c TopicsAuthorizationRuleClient) TopicsListKeys(ctx context.Context, id TopicAuthorizationRuleId) (result TopicsListKeysOperationResponse, err error) {
	req, err := c.preparerForTopicsListKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topicsauthorizationrule.TopicsAuthorizationRuleClient", "TopicsListKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "topicsauthorizationrule.TopicsAuthorizationRuleClient", "TopicsListKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTopicsListKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topicsauthorizationrule.TopicsAuthorizationRuleClient", "TopicsListKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTopicsListKeys prepares the TopicsListKeys request.
func (c TopicsAuthorizationRuleClient) preparerForTopicsListKeys(ctx context.Context, id TopicAuthorizationRuleId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTopicsListKeys handles the response to the TopicsListKeys request. The method always
// closes the http.Response Body.
func (c TopicsAuthorizationRuleClient) responderForTopicsListKeys(resp *http.Response) (result TopicsListKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
