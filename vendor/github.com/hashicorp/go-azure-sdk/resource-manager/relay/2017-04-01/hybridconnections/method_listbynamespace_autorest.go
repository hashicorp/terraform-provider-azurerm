package hybridconnections

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

type ListByNamespaceOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]HybridConnection

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByNamespaceOperationResponse, error)
}

type ListByNamespaceCompleteResult struct {
	Items []HybridConnection
}

func (r ListByNamespaceOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByNamespaceOperationResponse) LoadMore(ctx context.Context) (resp ListByNamespaceOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByNamespace ...
func (c HybridConnectionsClient) ListByNamespace(ctx context.Context, id NamespaceId) (resp ListByNamespaceOperationResponse, err error) {
	req, err := c.preparerForListByNamespace(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hybridconnections.HybridConnectionsClient", "ListByNamespace", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "hybridconnections.HybridConnectionsClient", "ListByNamespace", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByNamespace(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hybridconnections.HybridConnectionsClient", "ListByNamespace", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByNamespaceComplete retrieves all of the results into a single object
func (c HybridConnectionsClient) ListByNamespaceComplete(ctx context.Context, id NamespaceId) (ListByNamespaceCompleteResult, error) {
	return c.ListByNamespaceCompleteMatchingPredicate(ctx, id, HybridConnectionOperationPredicate{})
}

// ListByNamespaceCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c HybridConnectionsClient) ListByNamespaceCompleteMatchingPredicate(ctx context.Context, id NamespaceId, predicate HybridConnectionOperationPredicate) (resp ListByNamespaceCompleteResult, err error) {
	items := make([]HybridConnection, 0)

	page, err := c.ListByNamespace(ctx, id)
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

	out := ListByNamespaceCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByNamespace prepares the ListByNamespace request.
func (c HybridConnectionsClient) preparerForListByNamespace(ctx context.Context, id NamespaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/hybridConnections", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByNamespaceWithNextLink prepares the ListByNamespace request with the given nextLink token.
func (c HybridConnectionsClient) preparerForListByNamespaceWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByNamespace handles the response to the ListByNamespace request. The method always
// closes the http.Response Body.
func (c HybridConnectionsClient) responderForListByNamespace(resp *http.Response) (result ListByNamespaceOperationResponse, err error) {
	type page struct {
		Values   []HybridConnection `json:"value"`
		NextLink *string            `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByNamespaceOperationResponse, err error) {
			req, err := c.preparerForListByNamespaceWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "hybridconnections.HybridConnectionsClient", "ListByNamespace", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "hybridconnections.HybridConnectionsClient", "ListByNamespace", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByNamespace(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "hybridconnections.HybridConnectionsClient", "ListByNamespace", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
