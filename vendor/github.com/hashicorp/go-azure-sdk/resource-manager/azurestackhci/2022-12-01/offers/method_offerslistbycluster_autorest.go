package offers

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

type OffersListByClusterOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Offer

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (OffersListByClusterOperationResponse, error)
}

type OffersListByClusterCompleteResult struct {
	Items []Offer
}

func (r OffersListByClusterOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r OffersListByClusterOperationResponse) LoadMore(ctx context.Context) (resp OffersListByClusterOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type OffersListByClusterOperationOptions struct {
	Expand *string
}

func DefaultOffersListByClusterOperationOptions() OffersListByClusterOperationOptions {
	return OffersListByClusterOperationOptions{}
}

func (o OffersListByClusterOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o OffersListByClusterOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Expand != nil {
		out["$expand"] = *o.Expand
	}

	return out
}

// OffersListByCluster ...
func (c OffersClient) OffersListByCluster(ctx context.Context, id ClusterId, options OffersListByClusterOperationOptions) (resp OffersListByClusterOperationResponse, err error) {
	req, err := c.preparerForOffersListByCluster(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "offers.OffersClient", "OffersListByCluster", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "offers.OffersClient", "OffersListByCluster", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForOffersListByCluster(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "offers.OffersClient", "OffersListByCluster", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForOffersListByCluster prepares the OffersListByCluster request.
func (c OffersClient) preparerForOffersListByCluster(ctx context.Context, id ClusterId, options OffersListByClusterOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/offers", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForOffersListByClusterWithNextLink prepares the OffersListByCluster request with the given nextLink token.
func (c OffersClient) preparerForOffersListByClusterWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForOffersListByCluster handles the response to the OffersListByCluster request. The method always
// closes the http.Response Body.
func (c OffersClient) responderForOffersListByCluster(resp *http.Response) (result OffersListByClusterOperationResponse, err error) {
	type page struct {
		Values   []Offer `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result OffersListByClusterOperationResponse, err error) {
			req, err := c.preparerForOffersListByClusterWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "offers.OffersClient", "OffersListByCluster", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "offers.OffersClient", "OffersListByCluster", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForOffersListByCluster(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "offers.OffersClient", "OffersListByCluster", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// OffersListByClusterComplete retrieves all of the results into a single object
func (c OffersClient) OffersListByClusterComplete(ctx context.Context, id ClusterId, options OffersListByClusterOperationOptions) (OffersListByClusterCompleteResult, error) {
	return c.OffersListByClusterCompleteMatchingPredicate(ctx, id, options, OfferOperationPredicate{})
}

// OffersListByClusterCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c OffersClient) OffersListByClusterCompleteMatchingPredicate(ctx context.Context, id ClusterId, options OffersListByClusterOperationOptions, predicate OfferOperationPredicate) (resp OffersListByClusterCompleteResult, err error) {
	items := make([]Offer, 0)

	page, err := c.OffersListByCluster(ctx, id, options)
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

	out := OffersListByClusterCompleteResult{
		Items: items,
	}
	return out, nil
}
