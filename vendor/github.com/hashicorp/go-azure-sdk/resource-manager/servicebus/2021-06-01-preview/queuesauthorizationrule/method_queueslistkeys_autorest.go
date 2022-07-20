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

type QueuesListKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *AccessKeys
}

// QueuesListKeys ...
func (c QueuesAuthorizationRuleClient) QueuesListKeys(ctx context.Context, id QueueAuthorizationRuleId) (result QueuesListKeysOperationResponse, err error) {
	req, err := c.preparerForQueuesListKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesListKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesListKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueuesListKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesListKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueuesListKeys prepares the QueuesListKeys request.
func (c QueuesAuthorizationRuleClient) preparerForQueuesListKeys(ctx context.Context, id QueueAuthorizationRuleId) (*http.Request, error) {
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

// responderForQueuesListKeys handles the response to the QueuesListKeys request. The method always
// closes the http.Response Body.
func (c QueuesAuthorizationRuleClient) responderForQueuesListKeys(resp *http.Response) (result QueuesListKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
