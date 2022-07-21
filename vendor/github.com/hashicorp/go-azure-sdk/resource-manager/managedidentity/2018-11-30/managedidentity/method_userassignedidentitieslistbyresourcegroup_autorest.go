package managedidentity

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

type UserAssignedIdentitiesListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Identity

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (UserAssignedIdentitiesListByResourceGroupOperationResponse, error)
}

type UserAssignedIdentitiesListByResourceGroupCompleteResult struct {
	Items []Identity
}

func (r UserAssignedIdentitiesListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r UserAssignedIdentitiesListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp UserAssignedIdentitiesListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// UserAssignedIdentitiesListByResourceGroup ...
func (c ManagedIdentityClient) UserAssignedIdentitiesListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp UserAssignedIdentitiesListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForUserAssignedIdentitiesListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForUserAssignedIdentitiesListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// UserAssignedIdentitiesListByResourceGroupComplete retrieves all of the results into a single object
func (c ManagedIdentityClient) UserAssignedIdentitiesListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (UserAssignedIdentitiesListByResourceGroupCompleteResult, error) {
	return c.UserAssignedIdentitiesListByResourceGroupCompleteMatchingPredicate(ctx, id, IdentityOperationPredicate{})
}

// UserAssignedIdentitiesListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ManagedIdentityClient) UserAssignedIdentitiesListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate IdentityOperationPredicate) (resp UserAssignedIdentitiesListByResourceGroupCompleteResult, err error) {
	items := make([]Identity, 0)

	page, err := c.UserAssignedIdentitiesListByResourceGroup(ctx, id)
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

	out := UserAssignedIdentitiesListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForUserAssignedIdentitiesListByResourceGroup prepares the UserAssignedIdentitiesListByResourceGroup request.
func (c ManagedIdentityClient) preparerForUserAssignedIdentitiesListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForUserAssignedIdentitiesListByResourceGroupWithNextLink prepares the UserAssignedIdentitiesListByResourceGroup request with the given nextLink token.
func (c ManagedIdentityClient) preparerForUserAssignedIdentitiesListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForUserAssignedIdentitiesListByResourceGroup handles the response to the UserAssignedIdentitiesListByResourceGroup request. The method always
// closes the http.Response Body.
func (c ManagedIdentityClient) responderForUserAssignedIdentitiesListByResourceGroup(resp *http.Response) (result UserAssignedIdentitiesListByResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []Identity `json:"value"`
		NextLink *string    `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result UserAssignedIdentitiesListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForUserAssignedIdentitiesListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForUserAssignedIdentitiesListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
