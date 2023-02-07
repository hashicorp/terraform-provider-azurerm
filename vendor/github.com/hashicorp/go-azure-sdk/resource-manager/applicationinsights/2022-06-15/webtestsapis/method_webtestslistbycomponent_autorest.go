package webtestsapis

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

type WebTestsListByComponentOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]WebTest

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (WebTestsListByComponentOperationResponse, error)
}

type WebTestsListByComponentCompleteResult struct {
	Items []WebTest
}

func (r WebTestsListByComponentOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r WebTestsListByComponentOperationResponse) LoadMore(ctx context.Context) (resp WebTestsListByComponentOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// WebTestsListByComponent ...
func (c WebTestsAPIsClient) WebTestsListByComponent(ctx context.Context, id ComponentId) (resp WebTestsListByComponentOperationResponse, err error) {
	req, err := c.preparerForWebTestsListByComponent(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsListByComponent", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsListByComponent", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForWebTestsListByComponent(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsListByComponent", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForWebTestsListByComponent prepares the WebTestsListByComponent request.
func (c WebTestsAPIsClient) preparerForWebTestsListByComponent(ctx context.Context, id ComponentId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/webTests", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForWebTestsListByComponentWithNextLink prepares the WebTestsListByComponent request with the given nextLink token.
func (c WebTestsAPIsClient) preparerForWebTestsListByComponentWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForWebTestsListByComponent handles the response to the WebTestsListByComponent request. The method always
// closes the http.Response Body.
func (c WebTestsAPIsClient) responderForWebTestsListByComponent(resp *http.Response) (result WebTestsListByComponentOperationResponse, err error) {
	type page struct {
		Values   []WebTest `json:"value"`
		NextLink *string   `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result WebTestsListByComponentOperationResponse, err error) {
			req, err := c.preparerForWebTestsListByComponentWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsListByComponent", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsListByComponent", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForWebTestsListByComponent(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "webtestsapis.WebTestsAPIsClient", "WebTestsListByComponent", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// WebTestsListByComponentComplete retrieves all of the results into a single object
func (c WebTestsAPIsClient) WebTestsListByComponentComplete(ctx context.Context, id ComponentId) (WebTestsListByComponentCompleteResult, error) {
	return c.WebTestsListByComponentCompleteMatchingPredicate(ctx, id, WebTestOperationPredicate{})
}

// WebTestsListByComponentCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c WebTestsAPIsClient) WebTestsListByComponentCompleteMatchingPredicate(ctx context.Context, id ComponentId, predicate WebTestOperationPredicate) (resp WebTestsListByComponentCompleteResult, err error) {
	items := make([]WebTest, 0)

	page, err := c.WebTestsListByComponent(ctx, id)
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

	out := WebTestsListByComponentCompleteResult{
		Items: items,
	}
	return out, nil
}
