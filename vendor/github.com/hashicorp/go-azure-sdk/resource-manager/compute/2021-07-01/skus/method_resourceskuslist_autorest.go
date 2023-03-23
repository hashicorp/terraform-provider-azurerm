package skus

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

type ResourceSkusListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ResourceSku

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ResourceSkusListOperationResponse, error)
}

type ResourceSkusListCompleteResult struct {
	Items []ResourceSku
}

func (r ResourceSkusListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ResourceSkusListOperationResponse) LoadMore(ctx context.Context) (resp ResourceSkusListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ResourceSkusListOperationOptions struct {
	Filter                   *string
	IncludeExtendedLocations *string
}

func DefaultResourceSkusListOperationOptions() ResourceSkusListOperationOptions {
	return ResourceSkusListOperationOptions{}
}

func (o ResourceSkusListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ResourceSkusListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.IncludeExtendedLocations != nil {
		out["includeExtendedLocations"] = *o.IncludeExtendedLocations
	}

	return out
}

// ResourceSkusList ...
func (c SkusClient) ResourceSkusList(ctx context.Context, id commonids.SubscriptionId, options ResourceSkusListOperationOptions) (resp ResourceSkusListOperationResponse, err error) {
	req, err := c.preparerForResourceSkusList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "skus.SkusClient", "ResourceSkusList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "skus.SkusClient", "ResourceSkusList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForResourceSkusList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "skus.SkusClient", "ResourceSkusList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForResourceSkusList prepares the ResourceSkusList request.
func (c SkusClient) preparerForResourceSkusList(ctx context.Context, id commonids.SubscriptionId, options ResourceSkusListOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Compute/skus", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForResourceSkusListWithNextLink prepares the ResourceSkusList request with the given nextLink token.
func (c SkusClient) preparerForResourceSkusListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForResourceSkusList handles the response to the ResourceSkusList request. The method always
// closes the http.Response Body.
func (c SkusClient) responderForResourceSkusList(resp *http.Response) (result ResourceSkusListOperationResponse, err error) {
	type page struct {
		Values   []ResourceSku `json:"value"`
		NextLink *string       `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ResourceSkusListOperationResponse, err error) {
			req, err := c.preparerForResourceSkusListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "skus.SkusClient", "ResourceSkusList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "skus.SkusClient", "ResourceSkusList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForResourceSkusList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "skus.SkusClient", "ResourceSkusList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ResourceSkusListComplete retrieves all of the results into a single object
func (c SkusClient) ResourceSkusListComplete(ctx context.Context, id commonids.SubscriptionId, options ResourceSkusListOperationOptions) (ResourceSkusListCompleteResult, error) {
	return c.ResourceSkusListCompleteMatchingPredicate(ctx, id, options, ResourceSkuOperationPredicate{})
}

// ResourceSkusListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c SkusClient) ResourceSkusListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options ResourceSkusListOperationOptions, predicate ResourceSkuOperationPredicate) (resp ResourceSkusListCompleteResult, err error) {
	items := make([]ResourceSku, 0)

	page, err := c.ResourceSkusList(ctx, id, options)
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

	out := ResourceSkusListCompleteResult{
		Items: items,
	}
	return out, nil
}
