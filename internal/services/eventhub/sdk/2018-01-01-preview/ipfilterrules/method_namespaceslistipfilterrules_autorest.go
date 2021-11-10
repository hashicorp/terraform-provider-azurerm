package ipfilterrules

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type NamespacesListIPFilterRulesResponse struct {
	HttpResponse *http.Response
	Model        *[]IpFilterRule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (NamespacesListIPFilterRulesResponse, error)
}

type NamespacesListIPFilterRulesCompleteResult struct {
	Items []IpFilterRule
}

func (r NamespacesListIPFilterRulesResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r NamespacesListIPFilterRulesResponse) LoadMore(ctx context.Context) (resp NamespacesListIPFilterRulesResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// NamespacesListIPFilterRules ...
func (c IpFilterRulesClient) NamespacesListIPFilterRules(ctx context.Context, id NamespaceId) (resp NamespacesListIPFilterRulesResponse, err error) {
	req, err := c.preparerForNamespacesListIPFilterRules(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "ipfilterrules.IpFilterRulesClient", "NamespacesListIPFilterRules", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "ipfilterrules.IpFilterRulesClient", "NamespacesListIPFilterRules", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForNamespacesListIPFilterRules(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "ipfilterrules.IpFilterRulesClient", "NamespacesListIPFilterRules", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// NamespacesListIPFilterRulesComplete retrieves all of the results into a single object
func (c IpFilterRulesClient) NamespacesListIPFilterRulesComplete(ctx context.Context, id NamespaceId) (NamespacesListIPFilterRulesCompleteResult, error) {
	return c.NamespacesListIPFilterRulesCompleteMatchingPredicate(ctx, id, IpFilterRulePredicate{})
}

// NamespacesListIPFilterRulesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c IpFilterRulesClient) NamespacesListIPFilterRulesCompleteMatchingPredicate(ctx context.Context, id NamespaceId, predicate IpFilterRulePredicate) (resp NamespacesListIPFilterRulesCompleteResult, err error) {
	items := make([]IpFilterRule, 0)

	page, err := c.NamespacesListIPFilterRules(ctx, id)
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

	out := NamespacesListIPFilterRulesCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForNamespacesListIPFilterRules prepares the NamespacesListIPFilterRules request.
func (c IpFilterRulesClient) preparerForNamespacesListIPFilterRules(ctx context.Context, id NamespaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/ipfilterrules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForNamespacesListIPFilterRulesWithNextLink prepares the NamespacesListIPFilterRules request with the given nextLink token.
func (c IpFilterRulesClient) preparerForNamespacesListIPFilterRulesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForNamespacesListIPFilterRules handles the response to the NamespacesListIPFilterRules request. The method always
// closes the http.Response Body.
func (c IpFilterRulesClient) responderForNamespacesListIPFilterRules(resp *http.Response) (result NamespacesListIPFilterRulesResponse, err error) {
	type page struct {
		Values   []IpFilterRule `json:"value"`
		NextLink *string        `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result NamespacesListIPFilterRulesResponse, err error) {
			req, err := c.preparerForNamespacesListIPFilterRulesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "ipfilterrules.IpFilterRulesClient", "NamespacesListIPFilterRules", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "ipfilterrules.IpFilterRulesClient", "NamespacesListIPFilterRules", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForNamespacesListIPFilterRules(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "ipfilterrules.IpFilterRulesClient", "NamespacesListIPFilterRules", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
