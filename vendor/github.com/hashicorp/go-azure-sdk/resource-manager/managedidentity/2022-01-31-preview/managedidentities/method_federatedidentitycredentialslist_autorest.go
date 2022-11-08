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

type FederatedIdentityCredentialsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]FederatedIdentityCredential

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (FederatedIdentityCredentialsListOperationResponse, error)
}

type FederatedIdentityCredentialsListCompleteResult struct {
	Items []FederatedIdentityCredential
}

func (r FederatedIdentityCredentialsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r FederatedIdentityCredentialsListOperationResponse) LoadMore(ctx context.Context) (resp FederatedIdentityCredentialsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type FederatedIdentityCredentialsListOperationOptions struct {
	Top *int64
}

func DefaultFederatedIdentityCredentialsListOperationOptions() FederatedIdentityCredentialsListOperationOptions {
	return FederatedIdentityCredentialsListOperationOptions{}
}

func (o FederatedIdentityCredentialsListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o FederatedIdentityCredentialsListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// FederatedIdentityCredentialsList ...
func (c ManagedIdentitiesClient) FederatedIdentityCredentialsList(ctx context.Context, id commonids.UserAssignedIdentityId, options FederatedIdentityCredentialsListOperationOptions) (resp FederatedIdentityCredentialsListOperationResponse, err error) {
	req, err := c.preparerForFederatedIdentityCredentialsList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "FederatedIdentityCredentialsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "FederatedIdentityCredentialsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForFederatedIdentityCredentialsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "FederatedIdentityCredentialsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForFederatedIdentityCredentialsList prepares the FederatedIdentityCredentialsList request.
func (c ManagedIdentitiesClient) preparerForFederatedIdentityCredentialsList(ctx context.Context, id commonids.UserAssignedIdentityId, options FederatedIdentityCredentialsListOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/federatedIdentityCredentials", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForFederatedIdentityCredentialsListWithNextLink prepares the FederatedIdentityCredentialsList request with the given nextLink token.
func (c ManagedIdentitiesClient) preparerForFederatedIdentityCredentialsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForFederatedIdentityCredentialsList handles the response to the FederatedIdentityCredentialsList request. The method always
// closes the http.Response Body.
func (c ManagedIdentitiesClient) responderForFederatedIdentityCredentialsList(resp *http.Response) (result FederatedIdentityCredentialsListOperationResponse, err error) {
	type page struct {
		Values   []FederatedIdentityCredential `json:"value"`
		NextLink *string                       `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result FederatedIdentityCredentialsListOperationResponse, err error) {
			req, err := c.preparerForFederatedIdentityCredentialsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "FederatedIdentityCredentialsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "FederatedIdentityCredentialsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForFederatedIdentityCredentialsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedidentities.ManagedIdentitiesClient", "FederatedIdentityCredentialsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// FederatedIdentityCredentialsListComplete retrieves all of the results into a single object
func (c ManagedIdentitiesClient) FederatedIdentityCredentialsListComplete(ctx context.Context, id commonids.UserAssignedIdentityId, options FederatedIdentityCredentialsListOperationOptions) (FederatedIdentityCredentialsListCompleteResult, error) {
	return c.FederatedIdentityCredentialsListCompleteMatchingPredicate(ctx, id, options, FederatedIdentityCredentialOperationPredicate{})
}

// FederatedIdentityCredentialsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ManagedIdentitiesClient) FederatedIdentityCredentialsListCompleteMatchingPredicate(ctx context.Context, id commonids.UserAssignedIdentityId, options FederatedIdentityCredentialsListOperationOptions, predicate FederatedIdentityCredentialOperationPredicate) (resp FederatedIdentityCredentialsListCompleteResult, err error) {
	items := make([]FederatedIdentityCredential, 0)

	page, err := c.FederatedIdentityCredentialsList(ctx, id, options)
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

	out := FederatedIdentityCredentialsListCompleteResult{
		Items: items,
	}
	return out, nil
}
