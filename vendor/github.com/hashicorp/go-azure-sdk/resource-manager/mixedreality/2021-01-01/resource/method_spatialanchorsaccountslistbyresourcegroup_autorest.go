package resource

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

type SpatialAnchorsAccountsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]SpatialAnchorsAccount

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (SpatialAnchorsAccountsListByResourceGroupOperationResponse, error)
}

type SpatialAnchorsAccountsListByResourceGroupCompleteResult struct {
	Items []SpatialAnchorsAccount
}

func (r SpatialAnchorsAccountsListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r SpatialAnchorsAccountsListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp SpatialAnchorsAccountsListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// SpatialAnchorsAccountsListByResourceGroup ...
func (c ResourceClient) SpatialAnchorsAccountsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp SpatialAnchorsAccountsListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForSpatialAnchorsAccountsListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForSpatialAnchorsAccountsListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// SpatialAnchorsAccountsListByResourceGroupComplete retrieves all of the results into a single object
func (c ResourceClient) SpatialAnchorsAccountsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (SpatialAnchorsAccountsListByResourceGroupCompleteResult, error) {
	return c.SpatialAnchorsAccountsListByResourceGroupCompleteMatchingPredicate(ctx, id, SpatialAnchorsAccountOperationPredicate{})
}

// SpatialAnchorsAccountsListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ResourceClient) SpatialAnchorsAccountsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate SpatialAnchorsAccountOperationPredicate) (resp SpatialAnchorsAccountsListByResourceGroupCompleteResult, err error) {
	items := make([]SpatialAnchorsAccount, 0)

	page, err := c.SpatialAnchorsAccountsListByResourceGroup(ctx, id)
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

	out := SpatialAnchorsAccountsListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForSpatialAnchorsAccountsListByResourceGroup prepares the SpatialAnchorsAccountsListByResourceGroup request.
func (c ResourceClient) preparerForSpatialAnchorsAccountsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.MixedReality/spatialAnchorsAccounts", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForSpatialAnchorsAccountsListByResourceGroupWithNextLink prepares the SpatialAnchorsAccountsListByResourceGroup request with the given nextLink token.
func (c ResourceClient) preparerForSpatialAnchorsAccountsListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForSpatialAnchorsAccountsListByResourceGroup handles the response to the SpatialAnchorsAccountsListByResourceGroup request. The method always
// closes the http.Response Body.
func (c ResourceClient) responderForSpatialAnchorsAccountsListByResourceGroup(resp *http.Response) (result SpatialAnchorsAccountsListByResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []SpatialAnchorsAccount `json:"value"`
		NextLink *string                 `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result SpatialAnchorsAccountsListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForSpatialAnchorsAccountsListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForSpatialAnchorsAccountsListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
