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

type UserAssignedIdentitiesListAssociatedResourcesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AzureResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (UserAssignedIdentitiesListAssociatedResourcesOperationResponse, error)
}

type UserAssignedIdentitiesListAssociatedResourcesCompleteResult struct {
	Items []AzureResource
}

func (r UserAssignedIdentitiesListAssociatedResourcesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r UserAssignedIdentitiesListAssociatedResourcesOperationResponse) LoadMore(ctx context.Context) (resp UserAssignedIdentitiesListAssociatedResourcesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type UserAssignedIdentitiesListAssociatedResourcesOperationOptions struct {
	Filter  *string
	Orderby *string
	Skip    *int64
	Top     *int64
}

func DefaultUserAssignedIdentitiesListAssociatedResourcesOperationOptions() UserAssignedIdentitiesListAssociatedResourcesOperationOptions {
	return UserAssignedIdentitiesListAssociatedResourcesOperationOptions{}
}

func (o UserAssignedIdentitiesListAssociatedResourcesOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o UserAssignedIdentitiesListAssociatedResourcesOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Orderby != nil {
		out["$orderby"] = *o.Orderby
	}

	if o.Skip != nil {
		out["$skip"] = *o.Skip
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// UserAssignedIdentitiesListAssociatedResources ...
func (c ManagedIdentitiesClient) UserAssignedIdentitiesListAssociatedResources(ctx context.Context, id commonids.UserAssignedIdentityId, options UserAssignedIdentitiesListAssociatedResourcesOperationOptions) (resp UserAssignedIdentitiesListAssociatedResourcesOperationResponse, err error) {
	req, err := c.preparerForUserAssignedIdentitiesListAssociatedResources(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesListAssociatedResources", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesListAssociatedResources", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForUserAssignedIdentitiesListAssociatedResources(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesListAssociatedResources", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForUserAssignedIdentitiesListAssociatedResources prepares the UserAssignedIdentitiesListAssociatedResources request.
func (c ManagedIdentitiesClient) preparerForUserAssignedIdentitiesListAssociatedResources(ctx context.Context, id commonids.UserAssignedIdentityId, options UserAssignedIdentitiesListAssociatedResourcesOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/listAssociatedResources", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForUserAssignedIdentitiesListAssociatedResourcesWithNextLink prepares the UserAssignedIdentitiesListAssociatedResources request with the given nextLink token.
func (c ManagedIdentitiesClient) preparerForUserAssignedIdentitiesListAssociatedResourcesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUserAssignedIdentitiesListAssociatedResources handles the response to the UserAssignedIdentitiesListAssociatedResources request. The method always
// closes the http.Response Body.
func (c ManagedIdentitiesClient) responderForUserAssignedIdentitiesListAssociatedResources(resp *http.Response) (result UserAssignedIdentitiesListAssociatedResourcesOperationResponse, err error) {
	type page struct {
		Values   []AzureResource `json:"value"`
		NextLink *string         `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result UserAssignedIdentitiesListAssociatedResourcesOperationResponse, err error) {
			req, err := c.preparerForUserAssignedIdentitiesListAssociatedResourcesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesListAssociatedResources", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesListAssociatedResources", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForUserAssignedIdentitiesListAssociatedResources(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "UserAssignedIdentitiesListAssociatedResources", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// UserAssignedIdentitiesListAssociatedResourcesComplete retrieves all of the results into a single object
func (c ManagedIdentitiesClient) UserAssignedIdentitiesListAssociatedResourcesComplete(ctx context.Context, id commonids.UserAssignedIdentityId, options UserAssignedIdentitiesListAssociatedResourcesOperationOptions) (UserAssignedIdentitiesListAssociatedResourcesCompleteResult, error) {
	return c.UserAssignedIdentitiesListAssociatedResourcesCompleteMatchingPredicate(ctx, id, options, AzureResourceOperationPredicate{})
}

// UserAssignedIdentitiesListAssociatedResourcesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ManagedIdentitiesClient) UserAssignedIdentitiesListAssociatedResourcesCompleteMatchingPredicate(ctx context.Context, id commonids.UserAssignedIdentityId, options UserAssignedIdentitiesListAssociatedResourcesOperationOptions, predicate AzureResourceOperationPredicate) (resp UserAssignedIdentitiesListAssociatedResourcesCompleteResult, err error) {
	items := make([]AzureResource, 0)

	page, err := c.UserAssignedIdentitiesListAssociatedResources(ctx, id, options)
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

	out := UserAssignedIdentitiesListAssociatedResourcesCompleteResult{
		Items: items,
	}
	return out, nil
}
