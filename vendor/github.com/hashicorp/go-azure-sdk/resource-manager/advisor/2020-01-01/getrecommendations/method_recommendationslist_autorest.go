package getrecommendations

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

type RecommendationsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ResourceRecommendationBase

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (RecommendationsListOperationResponse, error)
}

type RecommendationsListCompleteResult struct {
	Items []ResourceRecommendationBase
}

func (r RecommendationsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r RecommendationsListOperationResponse) LoadMore(ctx context.Context) (resp RecommendationsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type RecommendationsListOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultRecommendationsListOperationOptions() RecommendationsListOperationOptions {
	return RecommendationsListOperationOptions{}
}

func (o RecommendationsListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o RecommendationsListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// RecommendationsList ...
func (c GetRecommendationsClient) RecommendationsList(ctx context.Context, id commonids.SubscriptionId, options RecommendationsListOperationOptions) (resp RecommendationsListOperationResponse, err error) {
	req, err := c.preparerForRecommendationsList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "getrecommendations.GetRecommendationsClient", "RecommendationsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "getrecommendations.GetRecommendationsClient", "RecommendationsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForRecommendationsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "getrecommendations.GetRecommendationsClient", "RecommendationsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForRecommendationsList prepares the RecommendationsList request.
func (c GetRecommendationsClient) preparerForRecommendationsList(ctx context.Context, id commonids.SubscriptionId, options RecommendationsListOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Advisor/recommendations", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForRecommendationsListWithNextLink prepares the RecommendationsList request with the given nextLink token.
func (c GetRecommendationsClient) preparerForRecommendationsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForRecommendationsList handles the response to the RecommendationsList request. The method always
// closes the http.Response Body.
func (c GetRecommendationsClient) responderForRecommendationsList(resp *http.Response) (result RecommendationsListOperationResponse, err error) {
	type page struct {
		Values   []ResourceRecommendationBase `json:"value"`
		NextLink *string                      `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result RecommendationsListOperationResponse, err error) {
			req, err := c.preparerForRecommendationsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "getrecommendations.GetRecommendationsClient", "RecommendationsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "getrecommendations.GetRecommendationsClient", "RecommendationsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForRecommendationsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "getrecommendations.GetRecommendationsClient", "RecommendationsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// RecommendationsListComplete retrieves all of the results into a single object
func (c GetRecommendationsClient) RecommendationsListComplete(ctx context.Context, id commonids.SubscriptionId, options RecommendationsListOperationOptions) (RecommendationsListCompleteResult, error) {
	return c.RecommendationsListCompleteMatchingPredicate(ctx, id, options, ResourceRecommendationBaseOperationPredicate{})
}

// RecommendationsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c GetRecommendationsClient) RecommendationsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options RecommendationsListOperationOptions, predicate ResourceRecommendationBaseOperationPredicate) (resp RecommendationsListCompleteResult, err error) {
	items := make([]ResourceRecommendationBase, 0)

	page, err := c.RecommendationsList(ctx, id, options)
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

	out := RecommendationsListCompleteResult{
		Items: items,
	}
	return out, nil
}
