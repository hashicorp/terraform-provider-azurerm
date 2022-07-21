package redisenterprise

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

type DatabasesListByClusterOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Database

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (DatabasesListByClusterOperationResponse, error)
}

type DatabasesListByClusterCompleteResult struct {
	Items []Database
}

func (r DatabasesListByClusterOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r DatabasesListByClusterOperationResponse) LoadMore(ctx context.Context) (resp DatabasesListByClusterOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// DatabasesListByCluster ...
func (c RedisEnterpriseClient) DatabasesListByCluster(ctx context.Context, id RedisEnterpriseId) (resp DatabasesListByClusterOperationResponse, err error) {
	req, err := c.preparerForDatabasesListByCluster(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesListByCluster", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesListByCluster", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForDatabasesListByCluster(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesListByCluster", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// DatabasesListByClusterComplete retrieves all of the results into a single object
func (c RedisEnterpriseClient) DatabasesListByClusterComplete(ctx context.Context, id RedisEnterpriseId) (DatabasesListByClusterCompleteResult, error) {
	return c.DatabasesListByClusterCompleteMatchingPredicate(ctx, id, DatabaseOperationPredicate{})
}

// DatabasesListByClusterCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c RedisEnterpriseClient) DatabasesListByClusterCompleteMatchingPredicate(ctx context.Context, id RedisEnterpriseId, predicate DatabaseOperationPredicate) (resp DatabasesListByClusterCompleteResult, err error) {
	items := make([]Database, 0)

	page, err := c.DatabasesListByCluster(ctx, id)
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

	out := DatabasesListByClusterCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForDatabasesListByCluster prepares the DatabasesListByCluster request.
func (c RedisEnterpriseClient) preparerForDatabasesListByCluster(ctx context.Context, id RedisEnterpriseId) (*http.Request, error) {
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

// preparerForDatabasesListByClusterWithNextLink prepares the DatabasesListByCluster request with the given nextLink token.
func (c RedisEnterpriseClient) preparerForDatabasesListByClusterWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForDatabasesListByCluster handles the response to the DatabasesListByCluster request. The method always
// closes the http.Response Body.
func (c RedisEnterpriseClient) responderForDatabasesListByCluster(resp *http.Response) (result DatabasesListByClusterOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result DatabasesListByClusterOperationResponse, err error) {
			req, err := c.preparerForDatabasesListByClusterWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesListByCluster", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesListByCluster", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForDatabasesListByCluster(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "redisenterprise.RedisEnterpriseClient", "DatabasesListByCluster", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
