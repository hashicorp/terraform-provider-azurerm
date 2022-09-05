package containerinstance

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

type ContainerGroupsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ContainerGroup

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ContainerGroupsListByResourceGroupOperationResponse, error)
}

type ContainerGroupsListByResourceGroupCompleteResult struct {
	Items []ContainerGroup
}

func (r ContainerGroupsListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ContainerGroupsListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp ContainerGroupsListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ContainerGroupsListByResourceGroup ...
func (c ContainerInstanceClient) ContainerGroupsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp ContainerGroupsListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForContainerGroupsListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForContainerGroupsListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForContainerGroupsListByResourceGroup prepares the ContainerGroupsListByResourceGroup request.
func (c ContainerInstanceClient) preparerForContainerGroupsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.ContainerInstance/containerGroups", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForContainerGroupsListByResourceGroupWithNextLink prepares the ContainerGroupsListByResourceGroup request with the given nextLink token.
func (c ContainerInstanceClient) preparerForContainerGroupsListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForContainerGroupsListByResourceGroup handles the response to the ContainerGroupsListByResourceGroup request. The method always
// closes the http.Response Body.
func (c ContainerInstanceClient) responderForContainerGroupsListByResourceGroup(resp *http.Response) (result ContainerGroupsListByResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []ContainerGroup `json:"value"`
		NextLink *string          `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ContainerGroupsListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForContainerGroupsListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForContainerGroupsListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ContainerGroupsListByResourceGroupComplete retrieves all of the results into a single object
func (c ContainerInstanceClient) ContainerGroupsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (ContainerGroupsListByResourceGroupCompleteResult, error) {
	return c.ContainerGroupsListByResourceGroupCompleteMatchingPredicate(ctx, id, ContainerGroupOperationPredicate{})
}

// ContainerGroupsListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ContainerInstanceClient) ContainerGroupsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ContainerGroupOperationPredicate) (resp ContainerGroupsListByResourceGroupCompleteResult, err error) {
	items := make([]ContainerGroup, 0)

	page, err := c.ContainerGroupsListByResourceGroup(ctx, id)
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

	out := ContainerGroupsListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}
