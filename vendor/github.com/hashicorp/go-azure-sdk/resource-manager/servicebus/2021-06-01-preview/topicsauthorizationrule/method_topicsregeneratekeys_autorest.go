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

type TopicsRegenerateKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *AccessKeys
}

// TopicsRegenerateKeys ...
func (c TopicsAuthorizationRuleClient) TopicsRegenerateKeys(ctx context.Context, id TopicAuthorizationRuleId, input RegenerateAccessKeyParameters) (result TopicsRegenerateKeysOperationResponse, err error) {
	req, err := c.preparerForTopicsRegenerateKeys(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topicsauthorizationrule.TopicsAuthorizationRuleClient", "TopicsRegenerateKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "topicsauthorizationrule.TopicsAuthorizationRuleClient", "TopicsRegenerateKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTopicsRegenerateKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topicsauthorizationrule.TopicsAuthorizationRuleClient", "TopicsRegenerateKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTopicsRegenerateKeys prepares the TopicsRegenerateKeys request.
func (c TopicsAuthorizationRuleClient) preparerForTopicsRegenerateKeys(ctx context.Context, id TopicAuthorizationRuleId, input RegenerateAccessKeyParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/regenerateKeys", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTopicsRegenerateKeys handles the response to the TopicsRegenerateKeys request. The method always
// closes the http.Response Body.
func (c TopicsAuthorizationRuleClient) responderForTopicsRegenerateKeys(resp *http.Response) (result TopicsRegenerateKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
