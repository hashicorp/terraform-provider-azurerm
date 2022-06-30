package rules

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagRulesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]MonitoringTagRules

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (TagRulesListOperationResponse, error)
}

type TagRulesListCompleteResult struct {
	Items []MonitoringTagRules
}

func (r TagRulesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r TagRulesListOperationResponse) LoadMore(ctx context.Context) (resp TagRulesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// TagRulesList ...
func (c RulesClient) TagRulesList(ctx context.Context, id MonitorId) (resp TagRulesListOperationResponse, err error) {
	req, err := c.preparerForTagRulesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "TagRulesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "TagRulesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForTagRulesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "TagRulesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// TagRulesListComplete retrieves all of the results into a single object
func (c RulesClient) TagRulesListComplete(ctx context.Context, id MonitorId) (TagRulesListCompleteResult, error) {
	return c.TagRulesListCompleteMatchingPredicate(ctx, id, MonitoringTagRulesOperationPredicate{})
}

// TagRulesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c RulesClient) TagRulesListCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate MonitoringTagRulesOperationPredicate) (resp TagRulesListCompleteResult, err error) {
	items := make([]MonitoringTagRules, 0)

	page, err := c.TagRulesList(ctx, id)
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

	out := TagRulesListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForTagRulesList prepares the TagRulesList request.
func (c RulesClient) preparerForTagRulesList(ctx context.Context, id MonitorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/tagRules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForTagRulesListWithNextLink prepares the TagRulesList request with the given nextLink token.
func (c RulesClient) preparerForTagRulesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForTagRulesList handles the response to the TagRulesList request. The method always
// closes the http.Response Body.
func (c RulesClient) responderForTagRulesList(resp *http.Response) (result TagRulesListOperationResponse, err error) {
	type page struct {
		Values   []MonitoringTagRules `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result TagRulesListOperationResponse, err error) {
			req, err := c.preparerForTagRulesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "rules.RulesClient", "TagRulesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "rules.RulesClient", "TagRulesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForTagRulesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "rules.RulesClient", "TagRulesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
