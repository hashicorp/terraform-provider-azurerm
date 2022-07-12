package dedicatedhsms

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

type DedicatedHsmListOutboundNetworkDependenciesEndpointsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]OutboundEnvironmentEndpoint

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (DedicatedHsmListOutboundNetworkDependenciesEndpointsOperationResponse, error)
}

type DedicatedHsmListOutboundNetworkDependenciesEndpointsCompleteResult struct {
	Items []OutboundEnvironmentEndpoint
}

func (r DedicatedHsmListOutboundNetworkDependenciesEndpointsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r DedicatedHsmListOutboundNetworkDependenciesEndpointsOperationResponse) LoadMore(ctx context.Context) (resp DedicatedHsmListOutboundNetworkDependenciesEndpointsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// DedicatedHsmListOutboundNetworkDependenciesEndpoints ...
func (c DedicatedHsmsClient) DedicatedHsmListOutboundNetworkDependenciesEndpoints(ctx context.Context, id DedicatedHSMId) (resp DedicatedHsmListOutboundNetworkDependenciesEndpointsOperationResponse, err error) {
	req, err := c.preparerForDedicatedHsmListOutboundNetworkDependenciesEndpoints(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListOutboundNetworkDependenciesEndpoints", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListOutboundNetworkDependenciesEndpoints", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForDedicatedHsmListOutboundNetworkDependenciesEndpoints(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListOutboundNetworkDependenciesEndpoints", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// DedicatedHsmListOutboundNetworkDependenciesEndpointsComplete retrieves all of the results into a single object
func (c DedicatedHsmsClient) DedicatedHsmListOutboundNetworkDependenciesEndpointsComplete(ctx context.Context, id DedicatedHSMId) (DedicatedHsmListOutboundNetworkDependenciesEndpointsCompleteResult, error) {
	return c.DedicatedHsmListOutboundNetworkDependenciesEndpointsCompleteMatchingPredicate(ctx, id, OutboundEnvironmentEndpointOperationPredicate{})
}

// DedicatedHsmListOutboundNetworkDependenciesEndpointsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c DedicatedHsmsClient) DedicatedHsmListOutboundNetworkDependenciesEndpointsCompleteMatchingPredicate(ctx context.Context, id DedicatedHSMId, predicate OutboundEnvironmentEndpointOperationPredicate) (resp DedicatedHsmListOutboundNetworkDependenciesEndpointsCompleteResult, err error) {
	items := make([]OutboundEnvironmentEndpoint, 0)

	page, err := c.DedicatedHsmListOutboundNetworkDependenciesEndpoints(ctx, id)
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

	out := DedicatedHsmListOutboundNetworkDependenciesEndpointsCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForDedicatedHsmListOutboundNetworkDependenciesEndpoints prepares the DedicatedHsmListOutboundNetworkDependenciesEndpoints request.
func (c DedicatedHsmsClient) preparerForDedicatedHsmListOutboundNetworkDependenciesEndpoints(ctx context.Context, id DedicatedHSMId) (*http.Request, error) {
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

// preparerForDedicatedHsmListOutboundNetworkDependenciesEndpointsWithNextLink prepares the DedicatedHsmListOutboundNetworkDependenciesEndpoints request with the given nextLink token.
func (c DedicatedHsmsClient) preparerForDedicatedHsmListOutboundNetworkDependenciesEndpointsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForDedicatedHsmListOutboundNetworkDependenciesEndpoints handles the response to the DedicatedHsmListOutboundNetworkDependenciesEndpoints request. The method always
// closes the http.Response Body.
func (c DedicatedHsmsClient) responderForDedicatedHsmListOutboundNetworkDependenciesEndpoints(resp *http.Response) (result DedicatedHsmListOutboundNetworkDependenciesEndpointsOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result DedicatedHsmListOutboundNetworkDependenciesEndpointsOperationResponse, err error) {
			req, err := c.preparerForDedicatedHsmListOutboundNetworkDependenciesEndpointsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListOutboundNetworkDependenciesEndpoints", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListOutboundNetworkDependenciesEndpoints", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForDedicatedHsmListOutboundNetworkDependenciesEndpoints(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListOutboundNetworkDependenciesEndpoints", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
