package alertrules

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRulesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AlertRule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (AlertRulesListOperationResponse, error)
}

type AlertRulesListCompleteResult struct {
	Items []AlertRule
}

func (r AlertRulesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r AlertRulesListOperationResponse) LoadMore(ctx context.Context) (resp AlertRulesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// AlertRulesList ...
func (c AlertRulesClient) AlertRulesList(ctx context.Context, id WorkspaceId) (resp AlertRulesListOperationResponse, err error) {
	req, err := c.preparerForAlertRulesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertrules.AlertRulesClient", "AlertRulesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertrules.AlertRulesClient", "AlertRulesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForAlertRulesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertrules.AlertRulesClient", "AlertRulesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForAlertRulesList prepares the AlertRulesList request.
func (c AlertRulesClient) preparerForAlertRulesList(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.SecurityInsights/alertRules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForAlertRulesListWithNextLink prepares the AlertRulesList request with the given nextLink token.
func (c AlertRulesClient) preparerForAlertRulesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForAlertRulesList handles the response to the AlertRulesList request. The method always
// closes the http.Response Body.
func (c AlertRulesClient) responderForAlertRulesList(resp *http.Response) (result AlertRulesListOperationResponse, err error) {
	type page struct {
		Values   []json.RawMessage `json:"value"`
		NextLink *string           `json:"nextLink"`
	}
	var respObj page
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&respObj),
		autorest.ByClosing())
	result.HttpResponse = resp
	temp := make([]AlertRule, 0)
	for i, v := range respObj.Values {
		val, err := unmarshalAlertRuleImplementation(v)
		if err != nil {
			err = fmt.Errorf("unmarshalling item %d for AlertRule (%q): %+v", i, v, err)
			return result, err
		}
		temp = append(temp, val)
	}
	result.Model = &temp
	result.nextLink = respObj.NextLink
	if respObj.NextLink != nil {
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result AlertRulesListOperationResponse, err error) {
			req, err := c.preparerForAlertRulesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "alertrules.AlertRulesClient", "AlertRulesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "alertrules.AlertRulesClient", "AlertRulesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForAlertRulesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "alertrules.AlertRulesClient", "AlertRulesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// AlertRulesListComplete retrieves all of the results into a single object
func (c AlertRulesClient) AlertRulesListComplete(ctx context.Context, id WorkspaceId) (AlertRulesListCompleteResult, error) {
	return c.AlertRulesListCompleteMatchingPredicate(ctx, id, AlertRuleOperationPredicate{})
}

// AlertRulesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AlertRulesClient) AlertRulesListCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, predicate AlertRuleOperationPredicate) (resp AlertRulesListCompleteResult, err error) {
	items := make([]AlertRule, 0)

	page, err := c.AlertRulesList(ctx, id)
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

	out := AlertRulesListCompleteResult{
		Items: items,
	}
	return out, nil
}
