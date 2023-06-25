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

type ListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *LogSearchRuleResourceCollection
}

type ListByResourceGroupOperationOptions struct {
	Filter *string
}

func DefaultListByResourceGroupOperationOptions() ListByResourceGroupOperationOptions {
	return ListByResourceGroupOperationOptions{}
}

func (o ListByResourceGroupOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListByResourceGroupOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// ListByResourceGroup ...
func (c ScheduledQueryRulesClient) ListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options ListByResourceGroupOperationOptions) (result ListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForListByResourceGroup(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledqueryrules.ScheduledQueryRulesClient", "ListByResourceGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledqueryrules.ScheduledQueryRulesClient", "ListByResourceGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListByResourceGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledqueryrules.ScheduledQueryRulesClient", "ListByResourceGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListByResourceGroup prepares the ListByResourceGroup request.
func (c ScheduledQueryRulesClient) preparerForListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options ListByResourceGroupOperationOptions) (*http.Request, error) {
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

// responderForListByResourceGroup handles the response to the ListByResourceGroup request. The method always
// closes the http.Response Body.
func (c ScheduledQueryRulesClient) responderForListByResourceGroup(resp *http.Response) (result ListByResourceGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
