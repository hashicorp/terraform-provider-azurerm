package scheduledqueryrules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *LogSearchRuleResourceCollection
}

type ListBySubscriptionOperationOptions struct {
	Filter *string
}

func DefaultListBySubscriptionOperationOptions() ListBySubscriptionOperationOptions {
	return ListBySubscriptionOperationOptions{}
}

func (o ListBySubscriptionOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListBySubscriptionOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// ListBySubscription ...
func (c ScheduledQueryRulesClient) ListBySubscription(ctx context.Context, id commonids.SubscriptionId, options ListBySubscriptionOperationOptions) (result ListBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForListBySubscription(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledqueryrules.ScheduledQueryRulesClient", "ListBySubscription", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledqueryrules.ScheduledQueryRulesClient", "ListBySubscription", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListBySubscription(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledqueryrules.ScheduledQueryRulesClient", "ListBySubscription", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListBySubscription prepares the ListBySubscription request.
func (c ScheduledQueryRulesClient) preparerForListBySubscription(ctx context.Context, id commonids.SubscriptionId, options ListBySubscriptionOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Insights/scheduledQueryRules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListBySubscription handles the response to the ListBySubscription request. The method always
// closes the http.Response Body.
func (c ScheduledQueryRulesClient) responderForListBySubscription(resp *http.Response) (result ListBySubscriptionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
