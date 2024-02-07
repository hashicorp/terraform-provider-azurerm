package managedclusters

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

type ListMeshUpgradeProfilesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]MeshUpgradeProfile

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListMeshUpgradeProfilesOperationResponse, error)
}

type ListMeshUpgradeProfilesCompleteResult struct {
	Items []MeshUpgradeProfile
}

func (r ListMeshUpgradeProfilesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListMeshUpgradeProfilesOperationResponse) LoadMore(ctx context.Context) (resp ListMeshUpgradeProfilesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListMeshUpgradeProfiles ...
func (c ManagedClustersClient) ListMeshUpgradeProfiles(ctx context.Context, id commonids.KubernetesClusterId) (resp ListMeshUpgradeProfilesOperationResponse, err error) {
	req, err := c.preparerForListMeshUpgradeProfiles(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListMeshUpgradeProfiles", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListMeshUpgradeProfiles", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListMeshUpgradeProfiles(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListMeshUpgradeProfiles", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListMeshUpgradeProfiles prepares the ListMeshUpgradeProfiles request.
func (c ManagedClustersClient) preparerForListMeshUpgradeProfiles(ctx context.Context, id commonids.KubernetesClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/meshUpgradeProfiles", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListMeshUpgradeProfilesWithNextLink prepares the ListMeshUpgradeProfiles request with the given nextLink token.
func (c ManagedClustersClient) preparerForListMeshUpgradeProfilesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListMeshUpgradeProfiles handles the response to the ListMeshUpgradeProfiles request. The method always
// closes the http.Response Body.
func (c ManagedClustersClient) responderForListMeshUpgradeProfiles(resp *http.Response) (result ListMeshUpgradeProfilesOperationResponse, err error) {
	type page struct {
		Values   []MeshUpgradeProfile `json:"value"`
		NextLink *string              `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListMeshUpgradeProfilesOperationResponse, err error) {
			req, err := c.preparerForListMeshUpgradeProfilesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListMeshUpgradeProfiles", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListMeshUpgradeProfiles", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListMeshUpgradeProfiles(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListMeshUpgradeProfiles", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListMeshUpgradeProfilesComplete retrieves all of the results into a single object
func (c ManagedClustersClient) ListMeshUpgradeProfilesComplete(ctx context.Context, id commonids.KubernetesClusterId) (ListMeshUpgradeProfilesCompleteResult, error) {
	return c.ListMeshUpgradeProfilesCompleteMatchingPredicate(ctx, id, MeshUpgradeProfileOperationPredicate{})
}

// ListMeshUpgradeProfilesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ManagedClustersClient) ListMeshUpgradeProfilesCompleteMatchingPredicate(ctx context.Context, id commonids.KubernetesClusterId, predicate MeshUpgradeProfileOperationPredicate) (resp ListMeshUpgradeProfilesCompleteResult, err error) {
	items := make([]MeshUpgradeProfile, 0)

	page, err := c.ListMeshUpgradeProfiles(ctx, id)
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

	out := ListMeshUpgradeProfilesCompleteResult{
		Items: items,
	}
	return out, nil
}
