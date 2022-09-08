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

type ContainerGroupsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ContainerGroup

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ContainerGroupsListOperationResponse, error)
}

type ContainerGroupsListCompleteResult struct {
	Items []ContainerGroup
}

func (r ContainerGroupsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ContainerGroupsListOperationResponse) LoadMore(ctx context.Context) (resp ContainerGroupsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ContainerGroupsList ...
func (c ContainerInstanceClient) ContainerGroupsList(ctx context.Context, id commonids.SubscriptionId) (resp ContainerGroupsListOperationResponse, err error) {
	req, err := c.preparerForContainerGroupsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForContainerGroupsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForContainerGroupsList prepares the ContainerGroupsList request.
func (c ContainerInstanceClient) preparerForContainerGroupsList(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
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

// preparerForContainerGroupsListWithNextLink prepares the ContainerGroupsList request with the given nextLink token.
func (c ContainerInstanceClient) preparerForContainerGroupsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForContainerGroupsList handles the response to the ContainerGroupsList request. The method always
// closes the http.Response Body.
func (c ContainerInstanceClient) responderForContainerGroupsList(resp *http.Response) (result ContainerGroupsListOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ContainerGroupsListOperationResponse, err error) {
			req, err := c.preparerForContainerGroupsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForContainerGroupsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "ContainerGroupsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ContainerGroupsListComplete retrieves all of the results into a single object
func (c ContainerInstanceClient) ContainerGroupsListComplete(ctx context.Context, id commonids.SubscriptionId) (ContainerGroupsListCompleteResult, error) {
	return c.ContainerGroupsListCompleteMatchingPredicate(ctx, id, ContainerGroupOperationPredicate{})
}

// ContainerGroupsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ContainerInstanceClient) ContainerGroupsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ContainerGroupOperationPredicate) (resp ContainerGroupsListCompleteResult, err error) {
	items := make([]ContainerGroup, 0)

	page, err := c.ContainerGroupsList(ctx, id)
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

	out := ContainerGroupsListCompleteResult{
		Items: items,
	}
	return out, nil
}
