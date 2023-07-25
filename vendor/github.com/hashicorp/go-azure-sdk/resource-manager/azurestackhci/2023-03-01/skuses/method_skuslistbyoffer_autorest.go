package skuses

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

type SkusListByOfferOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Sku

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (SkusListByOfferOperationResponse, error)
}

type SkusListByOfferCompleteResult struct {
	Items []Sku
}

func (r SkusListByOfferOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r SkusListByOfferOperationResponse) LoadMore(ctx context.Context) (resp SkusListByOfferOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type SkusListByOfferOperationOptions struct {
	Expand *string
}

func DefaultSkusListByOfferOperationOptions() SkusListByOfferOperationOptions {
	return SkusListByOfferOperationOptions{}
}

func (o SkusListByOfferOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o SkusListByOfferOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Expand != nil {
		out["$expand"] = *o.Expand
	}

	return out
}

// SkusListByOffer ...
func (c SkusesClient) SkusListByOffer(ctx context.Context, id OfferId, options SkusListByOfferOperationOptions) (resp SkusListByOfferOperationResponse, err error) {
	req, err := c.preparerForSkusListByOffer(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "skuses.SkusesClient", "SkusListByOffer", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "skuses.SkusesClient", "SkusListByOffer", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForSkusListByOffer(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "skuses.SkusesClient", "SkusListByOffer", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForSkusListByOffer prepares the SkusListByOffer request.
func (c SkusesClient) preparerForSkusListByOffer(ctx context.Context, id OfferId, options SkusListByOfferOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/skus", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForSkusListByOfferWithNextLink prepares the SkusListByOffer request with the given nextLink token.
func (c SkusesClient) preparerForSkusListByOfferWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForSkusListByOffer handles the response to the SkusListByOffer request. The method always
// closes the http.Response Body.
func (c SkusesClient) responderForSkusListByOffer(resp *http.Response) (result SkusListByOfferOperationResponse, err error) {
	type page struct {
		Values   []Sku   `json:"value"`
		NextLink *string `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result SkusListByOfferOperationResponse, err error) {
			req, err := c.preparerForSkusListByOfferWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "skuses.SkusesClient", "SkusListByOffer", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "skuses.SkusesClient", "SkusListByOffer", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForSkusListByOffer(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "skuses.SkusesClient", "SkusListByOffer", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// SkusListByOfferComplete retrieves all of the results into a single object
func (c SkusesClient) SkusListByOfferComplete(ctx context.Context, id OfferId, options SkusListByOfferOperationOptions) (SkusListByOfferCompleteResult, error) {
	return c.SkusListByOfferCompleteMatchingPredicate(ctx, id, options, SkuOperationPredicate{})
}

// SkusListByOfferCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c SkusesClient) SkusListByOfferCompleteMatchingPredicate(ctx context.Context, id OfferId, options SkusListByOfferOperationOptions, predicate SkuOperationPredicate) (resp SkusListByOfferCompleteResult, err error) {
	items := make([]Sku, 0)

	page, err := c.SkusListByOffer(ctx, id, options)
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

	out := SkusListByOfferCompleteResult{
		Items: items,
	}
	return out, nil
}
