package dnsforwardingrulesets

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

type ListByVirtualNetworkOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]VirtualNetworkDnsForwardingRuleset

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByVirtualNetworkOperationResponse, error)
}

type ListByVirtualNetworkCompleteResult struct {
	Items []VirtualNetworkDnsForwardingRuleset
}

func (r ListByVirtualNetworkOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByVirtualNetworkOperationResponse) LoadMore(ctx context.Context) (resp ListByVirtualNetworkOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByVirtualNetworkOperationOptions struct {
	Top *int64
}

func DefaultListByVirtualNetworkOperationOptions() ListByVirtualNetworkOperationOptions {
	return ListByVirtualNetworkOperationOptions{}
}

func (o ListByVirtualNetworkOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListByVirtualNetworkOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListByVirtualNetwork ...
func (c DnsForwardingRulesetsClient) ListByVirtualNetwork(ctx context.Context, id commonids.VirtualNetworkId, options ListByVirtualNetworkOperationOptions) (resp ListByVirtualNetworkOperationResponse, err error) {
	req, err := c.preparerForListByVirtualNetwork(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dnsforwardingrulesets.DnsForwardingRulesetsClient", "ListByVirtualNetwork", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "dnsforwardingrulesets.DnsForwardingRulesetsClient", "ListByVirtualNetwork", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByVirtualNetwork(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dnsforwardingrulesets.DnsForwardingRulesetsClient", "ListByVirtualNetwork", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByVirtualNetwork prepares the ListByVirtualNetwork request.
func (c DnsForwardingRulesetsClient) preparerForListByVirtualNetwork(ctx context.Context, id commonids.VirtualNetworkId, options ListByVirtualNetworkOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/listDnsForwardingRulesets", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByVirtualNetworkWithNextLink prepares the ListByVirtualNetwork request with the given nextLink token.
func (c DnsForwardingRulesetsClient) preparerForListByVirtualNetworkWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListByVirtualNetwork handles the response to the ListByVirtualNetwork request. The method always
// closes the http.Response Body.
func (c DnsForwardingRulesetsClient) responderForListByVirtualNetwork(resp *http.Response) (result ListByVirtualNetworkOperationResponse, err error) {
	type page struct {
		Values   []VirtualNetworkDnsForwardingRuleset `json:"value"`
		NextLink *string                              `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByVirtualNetworkOperationResponse, err error) {
			req, err := c.preparerForListByVirtualNetworkWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "dnsforwardingrulesets.DnsForwardingRulesetsClient", "ListByVirtualNetwork", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "dnsforwardingrulesets.DnsForwardingRulesetsClient", "ListByVirtualNetwork", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByVirtualNetwork(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "dnsforwardingrulesets.DnsForwardingRulesetsClient", "ListByVirtualNetwork", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByVirtualNetworkComplete retrieves all of the results into a single object
func (c DnsForwardingRulesetsClient) ListByVirtualNetworkComplete(ctx context.Context, id commonids.VirtualNetworkId, options ListByVirtualNetworkOperationOptions) (ListByVirtualNetworkCompleteResult, error) {
	return c.ListByVirtualNetworkCompleteMatchingPredicate(ctx, id, options, VirtualNetworkDnsForwardingRulesetOperationPredicate{})
}

// ListByVirtualNetworkCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c DnsForwardingRulesetsClient) ListByVirtualNetworkCompleteMatchingPredicate(ctx context.Context, id commonids.VirtualNetworkId, options ListByVirtualNetworkOperationOptions, predicate VirtualNetworkDnsForwardingRulesetOperationPredicate) (resp ListByVirtualNetworkCompleteResult, err error) {
	items := make([]VirtualNetworkDnsForwardingRuleset, 0)

	page, err := c.ListByVirtualNetwork(ctx, id, options)
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

	out := ListByVirtualNetworkCompleteResult{
		Items: items,
	}
	return out, nil
}
