package eventhubsclusters

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ClustersListByResourceGroupResponse struct {
	HttpResponse *http.Response
	Model        *[]Cluster

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ClustersListByResourceGroupResponse, error)
}

type ClustersListByResourceGroupCompleteResult struct {
	Items []Cluster
}

func (r ClustersListByResourceGroupResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ClustersListByResourceGroupResponse) LoadMore(ctx context.Context) (resp ClustersListByResourceGroupResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ClusterPredicate struct {
	// TODO: implement me
}

func (p ClusterPredicate) Matches(input Cluster) bool {
	// TODO: implement me
	// if p.Name != nil && input.Name != *p.Name {
	// 	return false
	// }

	return true
}

// ClustersListByResourceGroup ...
func (c EventHubsClustersClient) ClustersListByResourceGroup(ctx context.Context, id ResourceGroupId) (resp ClustersListByResourceGroupResponse, err error) {
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

// ClustersListByResourceGroupCompleteMatchingPredicate retrieves all of the results into a single object
func (c EventHubsClustersClient) ClustersListByResourceGroupComplete(ctx context.Context, id ResourceGroupId) (ClustersListByResourceGroupCompleteResult, error) {
	return c.ClustersListByResourceGroupCompleteMatchingPredicate(ctx, id, ClusterPredicate{})
}

// ClustersListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EventHubsClustersClient) ClustersListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id ResourceGroupId, predicate ClusterPredicate) (resp ClustersListByResourceGroupCompleteResult, err error) {
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
func (c EventHubsClustersClient) preparerForClustersListByResourceGroup(ctx context.Context, id ResourceGroupId) (*http.Request, error) {
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
func (c EventHubsClustersClient) responderForClustersListByResourceGroup(resp *http.Response) (result ClustersListByResourceGroupResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ClustersListByResourceGroupResponse, err error) {
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
