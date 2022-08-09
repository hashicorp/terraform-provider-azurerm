package eventhubsclusters

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

type ClustersListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Cluster

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ClustersListByResourceGroupOperationResponse, error)
}

type ClustersListByResourceGroupCompleteResult struct {
	Items []Cluster
}

func (r ClustersListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ClustersListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp ClustersListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ClustersListByResourceGroup ...
func (c EventHubsClustersClient) ClustersListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp ClustersListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForClustersListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclusters.EventHubsClustersClient", "ClustersListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclusters.EventHubsClustersClient", "ClustersListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForClustersListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclusters.EventHubsClustersClient", "ClustersListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ClustersListByResourceGroupComplete retrieves all of the results into a single object
func (c EventHubsClustersClient) ClustersListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (ClustersListByResourceGroupCompleteResult, error) {
	return c.ClustersListByResourceGroupCompleteMatchingPredicate(ctx, id, ClusterOperationPredicate{})
}

// ClustersListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EventHubsClustersClient) ClustersListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ClusterOperationPredicate) (resp ClustersListByResourceGroupCompleteResult, err error) {
	items := make([]Cluster, 0)

	page, err := c.ClustersListByResourceGroup(ctx, id)
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

	out := ClustersListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForClustersListByResourceGroup prepares the ClustersListByResourceGroup request.
func (c EventHubsClustersClient) preparerForClustersListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.EventHub/clusters", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForClustersListByResourceGroupWithNextLink prepares the ClustersListByResourceGroup request with the given nextLink token.
func (c EventHubsClustersClient) preparerForClustersListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForClustersListByResourceGroup handles the response to the ClustersListByResourceGroup request. The method always
// closes the http.Response Body.
func (c EventHubsClustersClient) responderForClustersListByResourceGroup(resp *http.Response) (result ClustersListByResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []Cluster `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ClustersListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForClustersListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventhubsclusters.EventHubsClustersClient", "ClustersListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventhubsclusters.EventHubsClustersClient", "ClustersListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForClustersListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventhubsclusters.EventHubsClustersClient", "ClustersListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
