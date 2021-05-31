package authorizationrulesnamespaces

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesListAuthorizationRulesResponse struct {
	HttpResponse *http.Response
	Model        *[]AuthorizationRule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (NamespacesListAuthorizationRulesResponse, error)
}

type NamespacesListAuthorizationRulesCompleteResult struct {
	Items []AuthorizationRule
}

func (r NamespacesListAuthorizationRulesResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r NamespacesListAuthorizationRulesResponse) LoadMore(ctx context.Context) (resp NamespacesListAuthorizationRulesResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type AuthorizationRulePredicate struct {
	// TODO: implement me
}

func (p AuthorizationRulePredicate) Matches(input AuthorizationRule) bool {
	// TODO: implement me
	// if p.Name != nil && input.Name != *p.Name {
	// 	return false
	// }

	return true
}

// NamespacesListAuthorizationRules ...
func (c AuthorizationRulesNamespacesClient) NamespacesListAuthorizationRules(ctx context.Context, id NamespaceId) (resp NamespacesListAuthorizationRulesResponse, err error) {
	req, err := c.preparerForNamespacesListAuthorizationRules(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesListAuthorizationRules", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesListAuthorizationRules", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForNamespacesListAuthorizationRules(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesListAuthorizationRules", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// NamespacesListAuthorizationRulesCompleteMatchingPredicate retrieves all of the results into a single object
func (c AuthorizationRulesNamespacesClient) NamespacesListAuthorizationRulesComplete(ctx context.Context, id NamespaceId) (NamespacesListAuthorizationRulesCompleteResult, error) {
	return c.NamespacesListAuthorizationRulesCompleteMatchingPredicate(ctx, id, AuthorizationRulePredicate{})
}

// NamespacesListAuthorizationRulesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AuthorizationRulesNamespacesClient) NamespacesListAuthorizationRulesCompleteMatchingPredicate(ctx context.Context, id NamespaceId, predicate AuthorizationRulePredicate) (resp NamespacesListAuthorizationRulesCompleteResult, err error) {
	items := make([]AuthorizationRule, 0)

	page, err := c.NamespacesListAuthorizationRules(ctx, id)
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

	out := NamespacesListAuthorizationRulesCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForNamespacesListAuthorizationRules prepares the NamespacesListAuthorizationRules request.
func (c AuthorizationRulesNamespacesClient) preparerForNamespacesListAuthorizationRules(ctx context.Context, id NamespaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/authorizationRules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForNamespacesListAuthorizationRulesWithNextLink prepares the NamespacesListAuthorizationRules request with the given nextLink token.
func (c AuthorizationRulesNamespacesClient) preparerForNamespacesListAuthorizationRulesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForNamespacesListAuthorizationRules handles the response to the NamespacesListAuthorizationRules request. The method always
// closes the http.Response Body.
func (c AuthorizationRulesNamespacesClient) responderForNamespacesListAuthorizationRules(resp *http.Response) (result NamespacesListAuthorizationRulesResponse, err error) {
	type page struct {
		Values   []AuthorizationRule `json:"value"`
		NextLink *string             `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result NamespacesListAuthorizationRulesResponse, err error) {
			req, err := c.preparerForNamespacesListAuthorizationRulesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesListAuthorizationRules", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesListAuthorizationRules", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForNamespacesListAuthorizationRules(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "authorizationrulesnamespaces.AuthorizationRulesNamespacesClient", "NamespacesListAuthorizationRules", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
