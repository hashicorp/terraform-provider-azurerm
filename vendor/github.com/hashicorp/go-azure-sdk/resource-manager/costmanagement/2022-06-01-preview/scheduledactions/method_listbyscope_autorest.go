package scheduledactions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByScopeOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ScheduledAction

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByScopeOperationResponse, error)
}

type ListByScopeCompleteResult struct {
	Items []ScheduledAction
}

func (r ListByScopeOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByScopeOperationResponse) LoadMore(ctx context.Context) (resp ListByScopeOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByScopeOperationOptions struct {
	Filter *string
}

func DefaultListByScopeOperationOptions() ListByScopeOperationOptions {
	return ListByScopeOperationOptions{}
}

func (o ListByScopeOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListByScopeOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// ListByScope ...
func (c ScheduledActionsClient) ListByScope(ctx context.Context, id commonids.ScopeId, options ListByScopeOperationOptions) (resp ListByScopeOperationResponse, err error) {
	req, err := c.preparerForListByScope(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledactions.ScheduledActionsClient", "ListByScope", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledactions.ScheduledActionsClient", "ListByScope", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByScope(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "scheduledactions.ScheduledActionsClient", "ListByScope", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByScope prepares the ListByScope request.
func (c ScheduledActionsClient) preparerForListByScope(ctx context.Context, id commonids.ScopeId, options ListByScopeOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.CostManagement/scheduledActions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByScopeWithNextLink prepares the ListByScope request with the given nextLink token.
func (c ScheduledActionsClient) preparerForListByScopeWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
	uri, err := url.Parse(nextLink)
	if err != nil {
		return nil, fmt.Errorf("parsing nextLink %q: %+v", nextLink, err)
	}
	queryParameters := map[string]interface{}{}
	for k, v := range uri.Query() {
		if len(v) == 0 {
			continue
		}
		val := v[0]
		val = autorest.Encode("query", val)
		queryParameters[k] = val
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListByScope handles the response to the ListByScope request. The method always
// closes the http.Response Body.
func (c ScheduledActionsClient) responderForListByScope(resp *http.Response) (result ListByScopeOperationResponse, err error) {
	type page struct {
		Values   []ScheduledAction `json:"value"`
		NextLink *string           `json:"nextLink"`
	}
	var respObj page
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&respObj),
		autorest.ByClosing())
	result.HttpResponse = resp
	result.Model = &respObj.Values
	result.nextLink = respObj.NextLink
	if respObj.NextLink != nil {
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByScopeOperationResponse, err error) {
			req, err := c.preparerForListByScopeWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "scheduledactions.ScheduledActionsClient", "ListByScope", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "scheduledactions.ScheduledActionsClient", "ListByScope", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByScope(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "scheduledactions.ScheduledActionsClient", "ListByScope", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByScopeComplete retrieves all of the results into a single object
func (c ScheduledActionsClient) ListByScopeComplete(ctx context.Context, id commonids.ScopeId, options ListByScopeOperationOptions) (ListByScopeCompleteResult, error) {
	return c.ListByScopeCompleteMatchingPredicate(ctx, id, options, ScheduledActionOperationPredicate{})
}

// ListByScopeCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ScheduledActionsClient) ListByScopeCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, options ListByScopeOperationOptions, predicate ScheduledActionOperationPredicate) (resp ListByScopeCompleteResult, err error) {
	items := make([]ScheduledAction, 0)

	page, err := c.ListByScope(ctx, id, options)
	if err != nil {
		err = fmt.Errorf("loading the initial page: %+v", err)
		return
	}
	if page.Model != nil {
		for _, v := range *page.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	for page.HasMore() {
		page, err = page.LoadMore(ctx)
		if err != nil {
			err = fmt.Errorf("loading the next page: %+v", err)
			return
		}

		if page.Model != nil {
			for _, v := range *page.Model {
				if predicate.Matches(v) {
					items = append(items, v)
				}
			}
		}
	}

	out := ListByScopeCompleteResult{
		Items: items,
	}
	return out, nil
}
