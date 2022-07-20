package queuesauthorizationrule

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueuesRegenerateKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *AccessKeys
}

// QueuesRegenerateKeys ...
func (c QueuesAuthorizationRuleClient) QueuesRegenerateKeys(ctx context.Context, id QueueAuthorizationRuleId, input RegenerateAccessKeyParameters) (result QueuesRegenerateKeysOperationResponse, err error) {
	req, err := c.preparerForQueuesRegenerateKeys(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesRegenerateKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesRegenerateKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueuesRegenerateKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesRegenerateKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueuesRegenerateKeys prepares the QueuesRegenerateKeys request.
func (c QueuesAuthorizationRuleClient) preparerForQueuesRegenerateKeys(ctx context.Context, id QueueAuthorizationRuleId, input RegenerateAccessKeyParameters) (*http.Request, error) {
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

// responderForQueuesRegenerateKeys handles the response to the QueuesRegenerateKeys request. The method always
// closes the http.Response Body.
func (c QueuesAuthorizationRuleClient) responderForQueuesRegenerateKeys(resp *http.Response) (result QueuesRegenerateKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
