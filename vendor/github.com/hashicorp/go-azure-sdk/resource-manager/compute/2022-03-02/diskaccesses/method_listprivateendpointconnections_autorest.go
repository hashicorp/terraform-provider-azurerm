package diskaccesses

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

type ListPrivateEndpointConnectionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]PrivateEndpointConnection

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListPrivateEndpointConnectionsOperationResponse, error)
}

type ListPrivateEndpointConnectionsCompleteResult struct {
	Items []PrivateEndpointConnection
}

func (r ListPrivateEndpointConnectionsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListPrivateEndpointConnectionsOperationResponse) LoadMore(ctx context.Context) (resp ListPrivateEndpointConnectionsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListPrivateEndpointConnections ...
func (c DiskAccessesClient) ListPrivateEndpointConnections(ctx context.Context, id DiskAccessId) (resp ListPrivateEndpointConnectionsOperationResponse, err error) {
	req, err := c.preparerForListPrivateEndpointConnections(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskaccesses.DiskAccessesClient", "ListPrivateEndpointConnections", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskaccesses.DiskAccessesClient", "ListPrivateEndpointConnections", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListPrivateEndpointConnections(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskaccesses.DiskAccessesClient", "ListPrivateEndpointConnections", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListPrivateEndpointConnections prepares the ListPrivateEndpointConnections request.
func (c DiskAccessesClient) preparerForListPrivateEndpointConnections(ctx context.Context, id DiskAccessId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/privateEndpointConnections", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListPrivateEndpointConnectionsWithNextLink prepares the ListPrivateEndpointConnections request with the given nextLink token.
func (c DiskAccessesClient) preparerForListPrivateEndpointConnectionsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListPrivateEndpointConnections handles the response to the ListPrivateEndpointConnections request. The method always
// closes the http.Response Body.
func (c DiskAccessesClient) responderForListPrivateEndpointConnections(resp *http.Response) (result ListPrivateEndpointConnectionsOperationResponse, err error) {
	type page struct {
		Values   []PrivateEndpointConnection `json:"value"`
		NextLink *string                     `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListPrivateEndpointConnectionsOperationResponse, err error) {
			req, err := c.preparerForListPrivateEndpointConnectionsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "diskaccesses.DiskAccessesClient", "ListPrivateEndpointConnections", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "diskaccesses.DiskAccessesClient", "ListPrivateEndpointConnections", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListPrivateEndpointConnections(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "diskaccesses.DiskAccessesClient", "ListPrivateEndpointConnections", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListPrivateEndpointConnectionsComplete retrieves all of the results into a single object
func (c DiskAccessesClient) ListPrivateEndpointConnectionsComplete(ctx context.Context, id DiskAccessId) (ListPrivateEndpointConnectionsCompleteResult, error) {
	return c.ListPrivateEndpointConnectionsCompleteMatchingPredicate(ctx, id, PrivateEndpointConnectionOperationPredicate{})
}

// ListPrivateEndpointConnectionsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c DiskAccessesClient) ListPrivateEndpointConnectionsCompleteMatchingPredicate(ctx context.Context, id DiskAccessId, predicate PrivateEndpointConnectionOperationPredicate) (resp ListPrivateEndpointConnectionsCompleteResult, err error) {
	items := make([]PrivateEndpointConnection, 0)

	page, err := c.ListPrivateEndpointConnections(ctx, id)
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

	out := ListPrivateEndpointConnectionsCompleteResult{
		Items: items,
	}
	return out, nil
}
