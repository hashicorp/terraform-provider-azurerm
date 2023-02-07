package redis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallRulesDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// FirewallRulesDelete ...
func (c RedisClient) FirewallRulesDelete(ctx context.Context, id FirewallRuleId) (result FirewallRulesDeleteOperationResponse, err error) {
	req, err := c.preparerForFirewallRulesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "FirewallRulesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "FirewallRulesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForFirewallRulesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "FirewallRulesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForFirewallRulesDelete prepares the FirewallRulesDelete request.
func (c RedisClient) preparerForFirewallRulesDelete(ctx context.Context, id FirewallRuleId) (*http.Request, error) {
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

// responderForFirewallRulesDelete handles the response to the FirewallRulesDelete request. The method always
// closes the http.Response Body.
func (c RedisClient) responderForFirewallRulesDelete(resp *http.Response) (result FirewallRulesDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
