package connectedclusters

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

type ConnectedClusterListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ConnectedCluster

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ConnectedClusterListBySubscriptionOperationResponse, error)
}

type ConnectedClusterListBySubscriptionCompleteResult struct {
	Items []ConnectedCluster
}

func (r ConnectedClusterListBySubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ConnectedClusterListBySubscriptionOperationResponse) LoadMore(ctx context.Context) (resp ConnectedClusterListBySubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ConnectedClusterListBySubscription ...
func (c ConnectedClustersClient) ConnectedClusterListBySubscription(ctx context.Context, id commonids.SubscriptionId) (resp ConnectedClusterListBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForConnectedClusterListBySubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectedclusters.ConnectedClustersClient", "ConnectedClusterListBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectedclusters.ConnectedClustersClient", "ConnectedClusterListBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForConnectedClusterListBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectedclusters.ConnectedClustersClient", "ConnectedClusterListBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForConnectedClusterListBySubscription prepares the ConnectedClusterListBySubscription request.
func (c ConnectedClustersClient) preparerForConnectedClusterListBySubscription(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Kubernetes/connectedClusters", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForConnectedClusterListBySubscriptionWithNextLink prepares the ConnectedClusterListBySubscription request with the given nextLink token.
func (c ConnectedClustersClient) preparerForConnectedClusterListBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForConnectedClusterListBySubscription handles the response to the ConnectedClusterListBySubscription request. The method always
// closes the http.Response Body.
func (c ConnectedClustersClient) responderForConnectedClusterListBySubscription(resp *http.Response) (result ConnectedClusterListBySubscriptionOperationResponse, err error) {
	type page struct {
		Values   []ConnectedCluster `json:"value"`
		NextLink *string            `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ConnectedClusterListBySubscriptionOperationResponse, err error) {
			req, err := c.preparerForConnectedClusterListBySubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "connectedclusters.ConnectedClustersClient", "ConnectedClusterListBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "connectedclusters.ConnectedClustersClient", "ConnectedClusterListBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForConnectedClusterListBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "connectedclusters.ConnectedClustersClient", "ConnectedClusterListBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ConnectedClusterListBySubscriptionComplete retrieves all of the results into a single object
func (c ConnectedClustersClient) ConnectedClusterListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (ConnectedClusterListBySubscriptionCompleteResult, error) {
	return c.ConnectedClusterListBySubscriptionCompleteMatchingPredicate(ctx, id, ConnectedClusterOperationPredicate{})
}

// ConnectedClusterListBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ConnectedClustersClient) ConnectedClusterListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ConnectedClusterOperationPredicate) (resp ConnectedClusterListBySubscriptionCompleteResult, err error) {
	items := make([]ConnectedCluster, 0)

	page, err := c.ConnectedClusterListBySubscription(ctx, id)
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

	out := ConnectedClusterListBySubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}
