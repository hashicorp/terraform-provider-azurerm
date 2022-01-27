package rules

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByRuleSetResponse struct {
	HttpResponse *http.Response
	Model        *[]Rule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByRuleSetResponse, error)
}

type ListByRuleSetCompleteResult struct {
	Items []Rule
}

func (r ListByRuleSetResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByRuleSetResponse) LoadMore(ctx context.Context) (resp ListByRuleSetResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByRuleSet ...
func (c RulesClient) ListByRuleSet(ctx context.Context, id RuleSetId) (resp ListByRuleSetResponse, err error) {
	req, err := c.preparerForListByRuleSet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "ListByRuleSet", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "ListByRuleSet", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByRuleSet(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "ListByRuleSet", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByRuleSetComplete retrieves all of the results into a single object
func (c RulesClient) ListByRuleSetComplete(ctx context.Context, id RuleSetId) (ListByRuleSetCompleteResult, error) {
	return c.ListByRuleSetCompleteMatchingPredicate(ctx, id, RulePredicate{})
}

// ListByRuleSetCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c RulesClient) ListByRuleSetCompleteMatchingPredicate(ctx context.Context, id RuleSetId, predicate RulePredicate) (resp ListByRuleSetCompleteResult, err error) {
	items := make([]Rule, 0)

	page, err := c.ListByRuleSet(ctx, id)
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

	out := ListByRuleSetCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByRuleSet prepares the ListByRuleSet request.
func (c RulesClient) preparerForListByRuleSet(ctx context.Context, id RuleSetId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/rules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByRuleSetWithNextLink prepares the ListByRuleSet request with the given nextLink token.
func (c RulesClient) preparerForListByRuleSetWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByRuleSet handles the response to the ListByRuleSet request. The method always
// closes the http.Response Body.
func (c RulesClient) responderForListByRuleSet(resp *http.Response) (result ListByRuleSetResponse, err error) {
	type page struct {
		Values   []Rule  `json:"value"`
		NextLink *string `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByRuleSetResponse, err error) {
			req, err := c.preparerForListByRuleSetWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "rules.RulesClient", "ListByRuleSet", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "rules.RulesClient", "ListByRuleSet", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByRuleSet(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "rules.RulesClient", "ListByRuleSet", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
