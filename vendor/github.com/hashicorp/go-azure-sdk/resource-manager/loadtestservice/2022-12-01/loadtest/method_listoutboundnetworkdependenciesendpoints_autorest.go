package loadtest

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

type ListOutboundNetworkDependenciesEndpointsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]OutboundEnvironmentEndpoint

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListOutboundNetworkDependenciesEndpointsOperationResponse, error)
}

type ListOutboundNetworkDependenciesEndpointsCompleteResult struct {
	Items []OutboundEnvironmentEndpoint
}

func (r ListOutboundNetworkDependenciesEndpointsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListOutboundNetworkDependenciesEndpointsOperationResponse) LoadMore(ctx context.Context) (resp ListOutboundNetworkDependenciesEndpointsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListOutboundNetworkDependenciesEndpoints ...
func (c LoadTestClient) ListOutboundNetworkDependenciesEndpoints(ctx context.Context, id LoadTestId) (resp ListOutboundNetworkDependenciesEndpointsOperationResponse, err error) {
	req, err := c.preparerForListOutboundNetworkDependenciesEndpoints(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "loadtest.LoadTestClient", "ListOutboundNetworkDependenciesEndpoints", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "loadtest.LoadTestClient", "ListOutboundNetworkDependenciesEndpoints", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListOutboundNetworkDependenciesEndpoints(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "loadtest.LoadTestClient", "ListOutboundNetworkDependenciesEndpoints", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListOutboundNetworkDependenciesEndpoints prepares the ListOutboundNetworkDependenciesEndpoints request.
func (c LoadTestClient) preparerForListOutboundNetworkDependenciesEndpoints(ctx context.Context, id LoadTestId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/outboundNetworkDependenciesEndpoints", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListOutboundNetworkDependenciesEndpointsWithNextLink prepares the ListOutboundNetworkDependenciesEndpoints request with the given nextLink token.
func (c LoadTestClient) preparerForListOutboundNetworkDependenciesEndpointsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListOutboundNetworkDependenciesEndpoints handles the response to the ListOutboundNetworkDependenciesEndpoints request. The method always
// closes the http.Response Body.
func (c LoadTestClient) responderForListOutboundNetworkDependenciesEndpoints(resp *http.Response) (result ListOutboundNetworkDependenciesEndpointsOperationResponse, err error) {
	type page struct {
		Values   []OutboundEnvironmentEndpoint `json:"value"`
		NextLink *string                       `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListOutboundNetworkDependenciesEndpointsOperationResponse, err error) {
			req, err := c.preparerForListOutboundNetworkDependenciesEndpointsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "loadtest.LoadTestClient", "ListOutboundNetworkDependenciesEndpoints", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "loadtest.LoadTestClient", "ListOutboundNetworkDependenciesEndpoints", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListOutboundNetworkDependenciesEndpoints(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "loadtest.LoadTestClient", "ListOutboundNetworkDependenciesEndpoints", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListOutboundNetworkDependenciesEndpointsComplete retrieves all of the results into a single object
func (c LoadTestClient) ListOutboundNetworkDependenciesEndpointsComplete(ctx context.Context, id LoadTestId) (ListOutboundNetworkDependenciesEndpointsCompleteResult, error) {
	return c.ListOutboundNetworkDependenciesEndpointsCompleteMatchingPredicate(ctx, id, OutboundEnvironmentEndpointOperationPredicate{})
}

// ListOutboundNetworkDependenciesEndpointsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c LoadTestClient) ListOutboundNetworkDependenciesEndpointsCompleteMatchingPredicate(ctx context.Context, id LoadTestId, predicate OutboundEnvironmentEndpointOperationPredicate) (resp ListOutboundNetworkDependenciesEndpointsCompleteResult, err error) {
	items := make([]OutboundEnvironmentEndpoint, 0)

	page, err := c.ListOutboundNetworkDependenciesEndpoints(ctx, id)
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

	out := ListOutboundNetworkDependenciesEndpointsCompleteResult{
		Items: items,
	}
	return out, nil
}
