package trustedaccess

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

type RoleBindingsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]TrustedAccessRoleBinding

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (RoleBindingsListOperationResponse, error)
}

type RoleBindingsListCompleteResult struct {
	Items []TrustedAccessRoleBinding
}

func (r RoleBindingsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r RoleBindingsListOperationResponse) LoadMore(ctx context.Context) (resp RoleBindingsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// RoleBindingsList ...
func (c TrustedAccessClient) RoleBindingsList(ctx context.Context, id commonids.KubernetesClusterId) (resp RoleBindingsListOperationResponse, err error) {
	req, err := c.preparerForRoleBindingsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "trustedaccess.TrustedAccessClient", "RoleBindingsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "trustedaccess.TrustedAccessClient", "RoleBindingsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForRoleBindingsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "trustedaccess.TrustedAccessClient", "RoleBindingsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForRoleBindingsList prepares the RoleBindingsList request.
func (c TrustedAccessClient) preparerForRoleBindingsList(ctx context.Context, id commonids.KubernetesClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/trustedAccessRoleBindings", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForRoleBindingsListWithNextLink prepares the RoleBindingsList request with the given nextLink token.
func (c TrustedAccessClient) preparerForRoleBindingsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForRoleBindingsList handles the response to the RoleBindingsList request. The method always
// closes the http.Response Body.
func (c TrustedAccessClient) responderForRoleBindingsList(resp *http.Response) (result RoleBindingsListOperationResponse, err error) {
	type page struct {
		Values   []TrustedAccessRoleBinding `json:"value"`
		NextLink *string                    `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result RoleBindingsListOperationResponse, err error) {
			req, err := c.preparerForRoleBindingsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "trustedaccess.TrustedAccessClient", "RoleBindingsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "trustedaccess.TrustedAccessClient", "RoleBindingsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForRoleBindingsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "trustedaccess.TrustedAccessClient", "RoleBindingsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// RoleBindingsListComplete retrieves all of the results into a single object
func (c TrustedAccessClient) RoleBindingsListComplete(ctx context.Context, id commonids.KubernetesClusterId) (RoleBindingsListCompleteResult, error) {
	return c.RoleBindingsListCompleteMatchingPredicate(ctx, id, TrustedAccessRoleBindingOperationPredicate{})
}

// RoleBindingsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c TrustedAccessClient) RoleBindingsListCompleteMatchingPredicate(ctx context.Context, id commonids.KubernetesClusterId, predicate TrustedAccessRoleBindingOperationPredicate) (resp RoleBindingsListCompleteResult, err error) {
	items := make([]TrustedAccessRoleBinding, 0)

	page, err := c.RoleBindingsList(ctx, id)
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

	out := RoleBindingsListCompleteResult{
		Items: items,
	}
	return out, nil
}
