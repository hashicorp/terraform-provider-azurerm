package redis

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallRulesCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *RedisFirewallRule
}

// FirewallRulesCreateOrUpdate ...
func (c RedisClient) FirewallRulesCreateOrUpdate(ctx context.Context, id FirewallRuleId, input RedisFirewallRule) (result FirewallRulesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForFirewallRulesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "FirewallRulesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "FirewallRulesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForFirewallRulesCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "FirewallRulesCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForFirewallRulesCreateOrUpdate prepares the FirewallRulesCreateOrUpdate request.
func (c RedisClient) preparerForFirewallRulesCreateOrUpdate(ctx context.Context, id FirewallRuleId, input RedisFirewallRule) (*http.Request, error) {
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

// responderForFirewallRulesCreateOrUpdate handles the response to the FirewallRulesCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c RedisClient) responderForFirewallRulesCreateOrUpdate(resp *http.Response) (result FirewallRulesCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
