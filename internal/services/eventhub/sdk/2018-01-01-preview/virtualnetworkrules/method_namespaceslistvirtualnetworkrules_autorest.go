package virtualnetworkrules

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesListVirtualNetworkRulesResponse struct {
	HttpResponse *http.Response
	Model        *[]VirtualNetworkRule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (NamespacesListVirtualNetworkRulesResponse, error)
}

type NamespacesListVirtualNetworkRulesCompleteResult struct {
	Items []VirtualNetworkRule
}

func (r NamespacesListVirtualNetworkRulesResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r NamespacesListVirtualNetworkRulesResponse) LoadMore(ctx context.Context) (resp NamespacesListVirtualNetworkRulesResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// NamespacesListVirtualNetworkRules ...
func (c VirtualNetworkRulesClient) NamespacesListVirtualNetworkRules(ctx context.Context, id NamespaceId) (resp NamespacesListVirtualNetworkRulesResponse, err error) {
	req, err := c.preparerForNamespacesListVirtualNetworkRules(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualnetworkrules.VirtualNetworkRulesClient", "NamespacesListVirtualNetworkRules", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualnetworkrules.VirtualNetworkRulesClient", "NamespacesListVirtualNetworkRules", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForNamespacesListVirtualNetworkRules(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualnetworkrules.VirtualNetworkRulesClient", "NamespacesListVirtualNetworkRules", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// NamespacesListVirtualNetworkRulesComplete retrieves all of the results into a single object
func (c VirtualNetworkRulesClient) NamespacesListVirtualNetworkRulesComplete(ctx context.Context, id NamespaceId) (NamespacesListVirtualNetworkRulesCompleteResult, error) {
	return c.NamespacesListVirtualNetworkRulesCompleteMatchingPredicate(ctx, id, VirtualNetworkRulePredicate{})
}

// NamespacesListVirtualNetworkRulesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c VirtualNetworkRulesClient) NamespacesListVirtualNetworkRulesCompleteMatchingPredicate(ctx context.Context, id NamespaceId, predicate VirtualNetworkRulePredicate) (resp NamespacesListVirtualNetworkRulesCompleteResult, err error) {
	items := make([]VirtualNetworkRule, 0)

	page, err := c.NamespacesListVirtualNetworkRules(ctx, id)
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

	out := NamespacesListVirtualNetworkRulesCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForNamespacesListVirtualNetworkRules prepares the NamespacesListVirtualNetworkRules request.
func (c VirtualNetworkRulesClient) preparerForNamespacesListVirtualNetworkRules(ctx context.Context, id NamespaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/virtualnetworkrules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForNamespacesListVirtualNetworkRulesWithNextLink prepares the NamespacesListVirtualNetworkRules request with the given nextLink token.
func (c VirtualNetworkRulesClient) preparerForNamespacesListVirtualNetworkRulesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForNamespacesListVirtualNetworkRules handles the response to the NamespacesListVirtualNetworkRules request. The method always
// closes the http.Response Body.
func (c VirtualNetworkRulesClient) responderForNamespacesListVirtualNetworkRules(resp *http.Response) (result NamespacesListVirtualNetworkRulesResponse, err error) {
	type page struct {
		Values   []VirtualNetworkRule `json:"value"`
		NextLink *string              `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result NamespacesListVirtualNetworkRulesResponse, err error) {
			req, err := c.preparerForNamespacesListVirtualNetworkRulesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "virtualnetworkrules.VirtualNetworkRulesClient", "NamespacesListVirtualNetworkRules", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "virtualnetworkrules.VirtualNetworkRulesClient", "NamespacesListVirtualNetworkRules", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForNamespacesListVirtualNetworkRules(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "virtualnetworkrules.VirtualNetworkRulesClient", "NamespacesListVirtualNetworkRules", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
