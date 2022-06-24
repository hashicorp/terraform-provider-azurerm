package resourceguards

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

type GetResourcesInResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ResourceGuardResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (GetResourcesInResourceGroupOperationResponse, error)
}

type GetResourcesInResourceGroupCompleteResult struct {
	Items []ResourceGuardResource
}

func (r GetResourcesInResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r GetResourcesInResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp GetResourcesInResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// GetResourcesInResourceGroup ...
func (c ResourceGuardsClient) GetResourcesInResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp GetResourcesInResourceGroupOperationResponse, err error) {
	req, err := c.preparerForGetResourcesInResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetResourcesInResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetResourcesInResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForGetResourcesInResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetResourcesInResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// GetResourcesInResourceGroupComplete retrieves all of the results into a single object
func (c ResourceGuardsClient) GetResourcesInResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (GetResourcesInResourceGroupCompleteResult, error) {
	return c.GetResourcesInResourceGroupCompleteMatchingPredicate(ctx, id, ResourceGuardResourceOperationPredicate{})
}

// GetResourcesInResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ResourceGuardsClient) GetResourcesInResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ResourceGuardResourceOperationPredicate) (resp GetResourcesInResourceGroupCompleteResult, err error) {
	items := make([]ResourceGuardResource, 0)

	page, err := c.GetResourcesInResourceGroup(ctx, id)
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

	out := GetResourcesInResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForGetResourcesInResourceGroup prepares the GetResourcesInResourceGroup request.
func (c ResourceGuardsClient) preparerForGetResourcesInResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.DataProtection/resourceGuards", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForGetResourcesInResourceGroupWithNextLink prepares the GetResourcesInResourceGroup request with the given nextLink token.
func (c ResourceGuardsClient) preparerForGetResourcesInResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForGetResourcesInResourceGroup handles the response to the GetResourcesInResourceGroup request. The method always
// closes the http.Response Body.
func (c ResourceGuardsClient) responderForGetResourcesInResourceGroup(resp *http.Response) (result GetResourcesInResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []ResourceGuardResource `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result GetResourcesInResourceGroupOperationResponse, err error) {
			req, err := c.preparerForGetResourcesInResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetResourcesInResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetResourcesInResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForGetResourcesInResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetResourcesInResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
