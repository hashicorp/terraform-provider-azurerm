package capacitypools

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

type PoolsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]CapacityPool

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (PoolsListOperationResponse, error)
}

type PoolsListCompleteResult struct {
	Items []CapacityPool
}

func (r PoolsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r PoolsListOperationResponse) LoadMore(ctx context.Context) (resp PoolsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// PoolsList ...
func (c CapacityPoolsClient) PoolsList(ctx context.Context, id NetAppAccountId) (resp PoolsListOperationResponse, err error) {
	req, err := c.preparerForPoolsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacitypools.CapacityPoolsClient", "PoolsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacitypools.CapacityPoolsClient", "PoolsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForPoolsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "capacitypools.CapacityPoolsClient", "PoolsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// PoolsListComplete retrieves all of the results into a single object
func (c CapacityPoolsClient) PoolsListComplete(ctx context.Context, id NetAppAccountId) (PoolsListCompleteResult, error) {
	return c.PoolsListCompleteMatchingPredicate(ctx, id, CapacityPoolOperationPredicate{})
}

// PoolsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c CapacityPoolsClient) PoolsListCompleteMatchingPredicate(ctx context.Context, id NetAppAccountId, predicate CapacityPoolOperationPredicate) (resp PoolsListCompleteResult, err error) {
	items := make([]CapacityPool, 0)

	page, err := c.PoolsList(ctx, id)
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

	out := PoolsListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForPoolsList prepares the PoolsList request.
func (c CapacityPoolsClient) preparerForPoolsList(ctx context.Context, id NetAppAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/capacityPools", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForPoolsListWithNextLink prepares the PoolsList request with the given nextLink token.
func (c CapacityPoolsClient) preparerForPoolsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForPoolsList handles the response to the PoolsList request. The method always
// closes the http.Response Body.
func (c CapacityPoolsClient) responderForPoolsList(resp *http.Response) (result PoolsListOperationResponse, err error) {
	type page struct {
		Values   []CapacityPool `json:"value"`
		NextLink *string        `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result PoolsListOperationResponse, err error) {
			req, err := c.preparerForPoolsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "capacitypools.CapacityPoolsClient", "PoolsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "capacitypools.CapacityPoolsClient", "PoolsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForPoolsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "capacitypools.CapacityPoolsClient", "PoolsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
