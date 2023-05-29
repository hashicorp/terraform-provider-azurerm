package replicationprotectableitems

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByReplicationProtectionContainersOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ProtectableItem

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByReplicationProtectionContainersOperationResponse, error)
}

type ListByReplicationProtectionContainersCompleteResult struct {
	Items []ProtectableItem
}

func (r ListByReplicationProtectionContainersOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByReplicationProtectionContainersOperationResponse) LoadMore(ctx context.Context) (resp ListByReplicationProtectionContainersOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByReplicationProtectionContainersOperationOptions struct {
	Filter *string
	Take   *string
}

func DefaultListByReplicationProtectionContainersOperationOptions() ListByReplicationProtectionContainersOperationOptions {
	return ListByReplicationProtectionContainersOperationOptions{}
}

func (o ListByReplicationProtectionContainersOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListByReplicationProtectionContainersOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Take != nil {
		out["$take"] = *o.Take
	}

	return out
}

// ListByReplicationProtectionContainers ...
func (c ReplicationProtectableItemsClient) ListByReplicationProtectionContainers(ctx context.Context, id ReplicationProtectionContainerId, options ListByReplicationProtectionContainersOperationOptions) (resp ListByReplicationProtectionContainersOperationResponse, err error) {
	req, err := c.preparerForListByReplicationProtectionContainers(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotectableitems.ReplicationProtectableItemsClient", "ListByReplicationProtectionContainers", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotectableitems.ReplicationProtectableItemsClient", "ListByReplicationProtectionContainers", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByReplicationProtectionContainers(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationprotectableitems.ReplicationProtectableItemsClient", "ListByReplicationProtectionContainers", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByReplicationProtectionContainers prepares the ListByReplicationProtectionContainers request.
func (c ReplicationProtectableItemsClient) preparerForListByReplicationProtectionContainers(ctx context.Context, id ReplicationProtectionContainerId, options ListByReplicationProtectionContainersOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/replicationProtectableItems", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByReplicationProtectionContainersWithNextLink prepares the ListByReplicationProtectionContainers request with the given nextLink token.
func (c ReplicationProtectableItemsClient) preparerForListByReplicationProtectionContainersWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByReplicationProtectionContainers handles the response to the ListByReplicationProtectionContainers request. The method always
// closes the http.Response Body.
func (c ReplicationProtectableItemsClient) responderForListByReplicationProtectionContainers(resp *http.Response) (result ListByReplicationProtectionContainersOperationResponse, err error) {
	type page struct {
		Values   []ProtectableItem `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByReplicationProtectionContainersOperationResponse, err error) {
			req, err := c.preparerForListByReplicationProtectionContainersWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "replicationprotectableitems.ReplicationProtectableItemsClient", "ListByReplicationProtectionContainers", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "replicationprotectableitems.ReplicationProtectableItemsClient", "ListByReplicationProtectionContainers", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByReplicationProtectionContainers(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "replicationprotectableitems.ReplicationProtectableItemsClient", "ListByReplicationProtectionContainers", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByReplicationProtectionContainersComplete retrieves all of the results into a single object
func (c ReplicationProtectableItemsClient) ListByReplicationProtectionContainersComplete(ctx context.Context, id ReplicationProtectionContainerId, options ListByReplicationProtectionContainersOperationOptions) (ListByReplicationProtectionContainersCompleteResult, error) {
	return c.ListByReplicationProtectionContainersCompleteMatchingPredicate(ctx, id, options, ProtectableItemOperationPredicate{})
}

// ListByReplicationProtectionContainersCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ReplicationProtectableItemsClient) ListByReplicationProtectionContainersCompleteMatchingPredicate(ctx context.Context, id ReplicationProtectionContainerId, options ListByReplicationProtectionContainersOperationOptions, predicate ProtectableItemOperationPredicate) (resp ListByReplicationProtectionContainersCompleteResult, err error) {
	items := make([]ProtectableItem, 0)

	page, err := c.ListByReplicationProtectionContainers(ctx, id, options)
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

	out := ListByReplicationProtectionContainersCompleteResult{
		Items: items,
	}
	return out, nil
}
