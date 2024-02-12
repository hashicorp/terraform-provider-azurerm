package managedclusters

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

type ListMeshRevisionProfilesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]MeshRevisionProfile

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListMeshRevisionProfilesOperationResponse, error)
}

type ListMeshRevisionProfilesCompleteResult struct {
	Items []MeshRevisionProfile
}

func (r ListMeshRevisionProfilesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListMeshRevisionProfilesOperationResponse) LoadMore(ctx context.Context) (resp ListMeshRevisionProfilesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListMeshRevisionProfiles ...
func (c ManagedClustersClient) ListMeshRevisionProfiles(ctx context.Context, id LocationId) (resp ListMeshRevisionProfilesOperationResponse, err error) {
	req, err := c.preparerForListMeshRevisionProfiles(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListMeshRevisionProfiles", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListMeshRevisionProfiles", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListMeshRevisionProfiles(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListMeshRevisionProfiles", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListMeshRevisionProfiles prepares the ListMeshRevisionProfiles request.
func (c ManagedClustersClient) preparerForListMeshRevisionProfiles(ctx context.Context, id LocationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/meshRevisionProfiles", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListMeshRevisionProfilesWithNextLink prepares the ListMeshRevisionProfiles request with the given nextLink token.
func (c ManagedClustersClient) preparerForListMeshRevisionProfilesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListMeshRevisionProfiles handles the response to the ListMeshRevisionProfiles request. The method always
// closes the http.Response Body.
func (c ManagedClustersClient) responderForListMeshRevisionProfiles(resp *http.Response) (result ListMeshRevisionProfilesOperationResponse, err error) {
	type page struct {
		Values   []MeshRevisionProfile `json:"value"`
		NextLink *string               `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListMeshRevisionProfilesOperationResponse, err error) {
			req, err := c.preparerForListMeshRevisionProfilesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListMeshRevisionProfiles", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListMeshRevisionProfiles", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListMeshRevisionProfiles(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedclusters.ManagedClustersClient", "ListMeshRevisionProfiles", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListMeshRevisionProfilesComplete retrieves all of the results into a single object
func (c ManagedClustersClient) ListMeshRevisionProfilesComplete(ctx context.Context, id LocationId) (ListMeshRevisionProfilesCompleteResult, error) {
	return c.ListMeshRevisionProfilesCompleteMatchingPredicate(ctx, id, MeshRevisionProfileOperationPredicate{})
}

// ListMeshRevisionProfilesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ManagedClustersClient) ListMeshRevisionProfilesCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate MeshRevisionProfileOperationPredicate) (resp ListMeshRevisionProfilesCompleteResult, err error) {
	items := make([]MeshRevisionProfile, 0)

	page, err := c.ListMeshRevisionProfiles(ctx, id)
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

	out := ListMeshRevisionProfilesCompleteResult{
		Items: items,
	}
	return out, nil
}
