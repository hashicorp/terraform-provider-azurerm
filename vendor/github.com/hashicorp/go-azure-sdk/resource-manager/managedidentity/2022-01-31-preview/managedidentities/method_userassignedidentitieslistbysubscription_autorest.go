package managedidentities

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

type UserAssignedIdentitiesListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Identity

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (UserAssignedIdentitiesListBySubscriptionOperationResponse, error)
}

type UserAssignedIdentitiesListBySubscriptionCompleteResult struct {
	Items []Identity
}

func (r UserAssignedIdentitiesListBySubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r UserAssignedIdentitiesListBySubscriptionOperationResponse) LoadMore(ctx context.Context) (resp UserAssignedIdentitiesListBySubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// UserAssignedIdentitiesListBySubscription ...
func (c ManagedIdentitiesClient) UserAssignedIdentitiesListBySubscription(ctx context.Context, id commonids.SubscriptionId) (resp UserAssignedIdentitiesListBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForUserAssignedIdentitiesListBySubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesListBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesListBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForUserAssignedIdentitiesListBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesListBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForUserAssignedIdentitiesListBySubscription prepares the UserAssignedIdentitiesListBySubscription request.
func (c ManagedIdentitiesClient) preparerForUserAssignedIdentitiesListBySubscription(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
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

// preparerForUserAssignedIdentitiesListBySubscriptionWithNextLink prepares the UserAssignedIdentitiesListBySubscription request with the given nextLink token.
func (c ManagedIdentitiesClient) preparerForUserAssignedIdentitiesListBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForUserAssignedIdentitiesListBySubscription handles the response to the UserAssignedIdentitiesListBySubscription request. The method always
// closes the http.Response Body.
func (c ManagedIdentitiesClient) responderForUserAssignedIdentitiesListBySubscription(resp *http.Response) (result UserAssignedIdentitiesListBySubscriptionOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result UserAssignedIdentitiesListBySubscriptionOperationResponse, err error) {
			req, err := c.preparerForUserAssignedIdentitiesListBySubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesListBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesListBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForUserAssignedIdentitiesListBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesListBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// UserAssignedIdentitiesListBySubscriptionComplete retrieves all of the results into a single object
func (c ManagedIdentitiesClient) UserAssignedIdentitiesListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (UserAssignedIdentitiesListBySubscriptionCompleteResult, error) {
	return c.UserAssignedIdentitiesListBySubscriptionCompleteMatchingPredicate(ctx, id, IdentityOperationPredicate{})
}

// UserAssignedIdentitiesListBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ManagedIdentitiesClient) UserAssignedIdentitiesListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate IdentityOperationPredicate) (resp UserAssignedIdentitiesListBySubscriptionCompleteResult, err error) {
	items := make([]Identity, 0)

	page, err := c.UserAssignedIdentitiesListBySubscription(ctx, id)
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

	out := UserAssignedIdentitiesListBySubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}
