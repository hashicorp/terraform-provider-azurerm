package publishers

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

type PublishersListByClusterOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Publisher

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (PublishersListByClusterOperationResponse, error)
}

type PublishersListByClusterCompleteResult struct {
	Items []Publisher
}

func (r PublishersListByClusterOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r PublishersListByClusterOperationResponse) LoadMore(ctx context.Context) (resp PublishersListByClusterOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// PublishersListByCluster ...
func (c PublishersClient) PublishersListByCluster(ctx context.Context, id ClusterId) (resp PublishersListByClusterOperationResponse, err error) {
	req, err := c.preparerForPublishersListByCluster(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "publishers.PublishersClient", "PublishersListByCluster", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "publishers.PublishersClient", "PublishersListByCluster", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForPublishersListByCluster(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "publishers.PublishersClient", "PublishersListByCluster", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForPublishersListByCluster prepares the PublishersListByCluster request.
func (c PublishersClient) preparerForPublishersListByCluster(ctx context.Context, id ClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/publishers", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForPublishersListByClusterWithNextLink prepares the PublishersListByCluster request with the given nextLink token.
func (c PublishersClient) preparerForPublishersListByClusterWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForPublishersListByCluster handles the response to the PublishersListByCluster request. The method always
// closes the http.Response Body.
func (c PublishersClient) responderForPublishersListByCluster(resp *http.Response) (result PublishersListByClusterOperationResponse, err error) {
	type page struct {
		Values   []Publisher `json:"value"`
		NextLink *string     `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result PublishersListByClusterOperationResponse, err error) {
			req, err := c.preparerForPublishersListByClusterWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "publishers.PublishersClient", "PublishersListByCluster", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "publishers.PublishersClient", "PublishersListByCluster", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForPublishersListByCluster(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "publishers.PublishersClient", "PublishersListByCluster", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// PublishersListByClusterComplete retrieves all of the results into a single object
func (c PublishersClient) PublishersListByClusterComplete(ctx context.Context, id ClusterId) (PublishersListByClusterCompleteResult, error) {
	return c.PublishersListByClusterCompleteMatchingPredicate(ctx, id, PublisherOperationPredicate{})
}

// PublishersListByClusterCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c PublishersClient) PublishersListByClusterCompleteMatchingPredicate(ctx context.Context, id ClusterId, predicate PublisherOperationPredicate) (resp PublishersListByClusterCompleteResult, err error) {
	items := make([]Publisher, 0)

	page, err := c.PublishersListByCluster(ctx, id)
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

	out := PublishersListByClusterCompleteResult{
		Items: items,
	}
	return out, nil
}
