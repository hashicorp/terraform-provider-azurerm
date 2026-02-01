package webapplicationfirewallpolicies

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

type PoliciesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]WebApplicationFirewallPolicy

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (PoliciesListOperationResponse, error)
}

type PoliciesListCompleteResult struct {
	Items []WebApplicationFirewallPolicy
}

func (r PoliciesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r PoliciesListOperationResponse) LoadMore(ctx context.Context) (resp PoliciesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// PoliciesList ...
func (c WebApplicationFirewallPoliciesClient) PoliciesList(ctx context.Context, id commonids.ResourceGroupId) (resp PoliciesListOperationResponse, err error) {
	req, err := c.preparerForPoliciesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForPoliciesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForPoliciesList prepares the PoliciesList request.
func (c WebApplicationFirewallPoliciesClient) preparerForPoliciesList(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForPoliciesListWithNextLink prepares the PoliciesList request with the given nextLink token.
func (c WebApplicationFirewallPoliciesClient) preparerForPoliciesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForPoliciesList handles the response to the PoliciesList request. The method always
// closes the http.Response Body.
func (c WebApplicationFirewallPoliciesClient) responderForPoliciesList(resp *http.Response) (result PoliciesListOperationResponse, err error) {
	type page struct {
		Values   []WebApplicationFirewallPolicy `json:"value"`
		NextLink *string                        `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result PoliciesListOperationResponse, err error) {
			req, err := c.preparerForPoliciesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForPoliciesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// PoliciesListComplete retrieves all of the results into a single object
func (c WebApplicationFirewallPoliciesClient) PoliciesListComplete(ctx context.Context, id commonids.ResourceGroupId) (PoliciesListCompleteResult, error) {
	return c.PoliciesListCompleteMatchingPredicate(ctx, id, WebApplicationFirewallPolicyOperationPredicate{})
}

// PoliciesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c WebApplicationFirewallPoliciesClient) PoliciesListCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate WebApplicationFirewallPolicyOperationPredicate) (resp PoliciesListCompleteResult, err error) {
	items := make([]WebApplicationFirewallPolicy, 0)

	page, err := c.PoliciesList(ctx, id)
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

	out := PoliciesListCompleteResult{
		Items: items,
	}
	return out, nil
}
