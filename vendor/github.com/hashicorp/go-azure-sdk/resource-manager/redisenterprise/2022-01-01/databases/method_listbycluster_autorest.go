package databases

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

type ListByClusterOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Database

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByClusterOperationResponse, error)
}

type ListByClusterCompleteResult struct {
	Items []Database
}

func (r ListByClusterOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByClusterOperationResponse) LoadMore(ctx context.Context) (resp ListByClusterOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByCluster ...
func (c DatabasesClient) ListByCluster(ctx context.Context, id RedisEnterpriseId) (resp ListByClusterOperationResponse, err error) {
	req, err := c.preparerForListByCluster(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "ListByCluster", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "ListByCluster", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByCluster(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "ListByCluster", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByClusterComplete retrieves all of the results into a single object
func (c DatabasesClient) ListByClusterComplete(ctx context.Context, id RedisEnterpriseId) (ListByClusterCompleteResult, error) {
	return c.ListByClusterCompleteMatchingPredicate(ctx, id, DatabaseOperationPredicate{})
}

// ListByClusterCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c DatabasesClient) ListByClusterCompleteMatchingPredicate(ctx context.Context, id RedisEnterpriseId, predicate DatabaseOperationPredicate) (resp ListByClusterCompleteResult, err error) {
	items := make([]Database, 0)

	page, err := c.ListByCluster(ctx, id)
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

	out := ListByClusterCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByCluster prepares the ListByCluster request.
func (c DatabasesClient) preparerForListByCluster(ctx context.Context, id RedisEnterpriseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/databases", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByClusterWithNextLink prepares the ListByCluster request with the given nextLink token.
func (c DatabasesClient) preparerForListByClusterWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByCluster handles the response to the ListByCluster request. The method always
// closes the http.Response Body.
func (c DatabasesClient) responderForListByCluster(resp *http.Response) (result ListByClusterOperationResponse, err error) {
	type page struct {
		Values   []Database `json:"value"`
		NextLink *string    `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByClusterOperationResponse, err error) {
			req, err := c.preparerForListByClusterWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "ListByCluster", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "ListByCluster", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByCluster(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "ListByCluster", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
