package webapplicationfirewallmanagedrulesets

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ManagedRuleSetsListResponse struct {
	HttpResponse *http.Response
	Model        *[]ManagedRuleSetDefinition

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ManagedRuleSetsListResponse, error)
}

type ManagedRuleSetsListCompleteResult struct {
	Items []ManagedRuleSetDefinition
}

func (r ManagedRuleSetsListResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ManagedRuleSetsListResponse) LoadMore(ctx context.Context) (resp ManagedRuleSetsListResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ManagedRuleSetsList ...
func (c WebApplicationFirewallManagedRuleSetsClient) ManagedRuleSetsList(ctx context.Context, id SubscriptionId) (resp ManagedRuleSetsListResponse, err error) {
	req, err := c.preparerForManagedRuleSetsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallmanagedrulesets.WebApplicationFirewallManagedRuleSetsClient", "ManagedRuleSetsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallmanagedrulesets.WebApplicationFirewallManagedRuleSetsClient", "ManagedRuleSetsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForManagedRuleSetsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallmanagedrulesets.WebApplicationFirewallManagedRuleSetsClient", "ManagedRuleSetsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ManagedRuleSetsListComplete retrieves all of the results into a single object
func (c WebApplicationFirewallManagedRuleSetsClient) ManagedRuleSetsListComplete(ctx context.Context, id SubscriptionId) (ManagedRuleSetsListCompleteResult, error) {
	return c.ManagedRuleSetsListCompleteMatchingPredicate(ctx, id, ManagedRuleSetDefinitionPredicate{})
}

// ManagedRuleSetsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c WebApplicationFirewallManagedRuleSetsClient) ManagedRuleSetsListCompleteMatchingPredicate(ctx context.Context, id SubscriptionId, predicate ManagedRuleSetDefinitionPredicate) (resp ManagedRuleSetsListCompleteResult, err error) {
	items := make([]ManagedRuleSetDefinition, 0)

	page, err := c.ManagedRuleSetsList(ctx, id)
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

	out := ManagedRuleSetsListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForManagedRuleSetsList prepares the ManagedRuleSetsList request.
func (c WebApplicationFirewallManagedRuleSetsClient) preparerForManagedRuleSetsList(ctx context.Context, id SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Network/frontDoorWebApplicationFirewallManagedRuleSets", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForManagedRuleSetsListWithNextLink prepares the ManagedRuleSetsList request with the given nextLink token.
func (c WebApplicationFirewallManagedRuleSetsClient) preparerForManagedRuleSetsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForManagedRuleSetsList handles the response to the ManagedRuleSetsList request. The method always
// closes the http.Response Body.
func (c WebApplicationFirewallManagedRuleSetsClient) responderForManagedRuleSetsList(resp *http.Response) (result ManagedRuleSetsListResponse, err error) {
	type page struct {
		Values   []ManagedRuleSetDefinition `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ManagedRuleSetsListResponse, err error) {
			req, err := c.preparerForManagedRuleSetsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "webapplicationfirewallmanagedrulesets.WebApplicationFirewallManagedRuleSetsClient", "ManagedRuleSetsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "webapplicationfirewallmanagedrulesets.WebApplicationFirewallManagedRuleSetsClient", "ManagedRuleSetsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForManagedRuleSetsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "webapplicationfirewallmanagedrulesets.WebApplicationFirewallManagedRuleSetsClient", "ManagedRuleSetsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
