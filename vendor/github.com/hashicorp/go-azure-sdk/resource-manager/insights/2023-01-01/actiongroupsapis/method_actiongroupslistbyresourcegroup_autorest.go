package actiongroupsapis

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

type ActionGroupsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ActionGroupResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ActionGroupsListByResourceGroupOperationResponse, error)
}

type ActionGroupsListByResourceGroupCompleteResult struct {
	Items []ActionGroupResource
}

func (r ActionGroupsListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ActionGroupsListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp ActionGroupsListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ActionGroupsListByResourceGroup ...
func (c ActionGroupsAPIsClient) ActionGroupsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp ActionGroupsListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForActionGroupsListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForActionGroupsListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForActionGroupsListByResourceGroup prepares the ActionGroupsListByResourceGroup request.
func (c ActionGroupsAPIsClient) preparerForActionGroupsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Insights/actionGroups", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForActionGroupsListByResourceGroupWithNextLink prepares the ActionGroupsListByResourceGroup request with the given nextLink token.
func (c ActionGroupsAPIsClient) preparerForActionGroupsListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForActionGroupsListByResourceGroup handles the response to the ActionGroupsListByResourceGroup request. The method always
// closes the http.Response Body.
func (c ActionGroupsAPIsClient) responderForActionGroupsListByResourceGroup(resp *http.Response) (result ActionGroupsListByResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []ActionGroupResource `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ActionGroupsListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForActionGroupsListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForActionGroupsListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "actiongroupsapis.ActionGroupsAPIsClient", "ActionGroupsListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ActionGroupsListByResourceGroupComplete retrieves all of the results into a single object
func (c ActionGroupsAPIsClient) ActionGroupsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (ActionGroupsListByResourceGroupCompleteResult, error) {
	return c.ActionGroupsListByResourceGroupCompleteMatchingPredicate(ctx, id, ActionGroupResourceOperationPredicate{})
}

// ActionGroupsListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ActionGroupsAPIsClient) ActionGroupsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ActionGroupResourceOperationPredicate) (resp ActionGroupsListByResourceGroupCompleteResult, err error) {
	items := make([]ActionGroupResource, 0)

	page, err := c.ActionGroupsListByResourceGroup(ctx, id)
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

	out := ActionGroupsListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}
