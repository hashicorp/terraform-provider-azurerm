package replicationnetworkmappings

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

type ListByReplicationNetworksOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]NetworkMapping

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByReplicationNetworksOperationResponse, error)
}

type ListByReplicationNetworksCompleteResult struct {
	Items []NetworkMapping
}

func (r ListByReplicationNetworksOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByReplicationNetworksOperationResponse) LoadMore(ctx context.Context) (resp ListByReplicationNetworksOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByReplicationNetworks ...
func (c ReplicationNetworkMappingsClient) ListByReplicationNetworks(ctx context.Context, id ReplicationNetworkId) (resp ListByReplicationNetworksOperationResponse, err error) {
	req, err := c.preparerForListByReplicationNetworks(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationnetworkmappings.ReplicationNetworkMappingsClient", "ListByReplicationNetworks", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationnetworkmappings.ReplicationNetworkMappingsClient", "ListByReplicationNetworks", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByReplicationNetworks(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationnetworkmappings.ReplicationNetworkMappingsClient", "ListByReplicationNetworks", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByReplicationNetworks prepares the ListByReplicationNetworks request.
func (c ReplicationNetworkMappingsClient) preparerForListByReplicationNetworks(ctx context.Context, id ReplicationNetworkId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/replicationNetworkMappings", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByReplicationNetworksWithNextLink prepares the ListByReplicationNetworks request with the given nextLink token.
func (c ReplicationNetworkMappingsClient) preparerForListByReplicationNetworksWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByReplicationNetworks handles the response to the ListByReplicationNetworks request. The method always
// closes the http.Response Body.
func (c ReplicationNetworkMappingsClient) responderForListByReplicationNetworks(resp *http.Response) (result ListByReplicationNetworksOperationResponse, err error) {
	type page struct {
		Values   []NetworkMapping `json:"value"`
		NextLink *string          `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByReplicationNetworksOperationResponse, err error) {
			req, err := c.preparerForListByReplicationNetworksWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "replicationnetworkmappings.ReplicationNetworkMappingsClient", "ListByReplicationNetworks", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "replicationnetworkmappings.ReplicationNetworkMappingsClient", "ListByReplicationNetworks", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByReplicationNetworks(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "replicationnetworkmappings.ReplicationNetworkMappingsClient", "ListByReplicationNetworks", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByReplicationNetworksComplete retrieves all of the results into a single object
func (c ReplicationNetworkMappingsClient) ListByReplicationNetworksComplete(ctx context.Context, id ReplicationNetworkId) (ListByReplicationNetworksCompleteResult, error) {
	return c.ListByReplicationNetworksCompleteMatchingPredicate(ctx, id, NetworkMappingOperationPredicate{})
}

// ListByReplicationNetworksCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ReplicationNetworkMappingsClient) ListByReplicationNetworksCompleteMatchingPredicate(ctx context.Context, id ReplicationNetworkId, predicate NetworkMappingOperationPredicate) (resp ListByReplicationNetworksCompleteResult, err error) {
	items := make([]NetworkMapping, 0)

	page, err := c.ListByReplicationNetworks(ctx, id)
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

	out := ListByReplicationNetworksCompleteResult{
		Items: items,
	}
	return out, nil
}
