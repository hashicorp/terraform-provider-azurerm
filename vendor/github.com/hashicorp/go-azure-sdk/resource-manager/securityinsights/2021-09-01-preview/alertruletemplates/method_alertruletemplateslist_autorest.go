package alertruletemplates

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

type AlertRuleTemplatesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AlertRuleTemplate

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (AlertRuleTemplatesListOperationResponse, error)
}

type AlertRuleTemplatesListCompleteResult struct {
	Items []AlertRuleTemplate
}

func (r AlertRuleTemplatesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r AlertRuleTemplatesListOperationResponse) LoadMore(ctx context.Context) (resp AlertRuleTemplatesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// AlertRuleTemplatesList ...
func (c AlertRuleTemplatesClient) AlertRuleTemplatesList(ctx context.Context, id WorkspaceId) (resp AlertRuleTemplatesListOperationResponse, err error) {
	req, err := c.preparerForAlertRuleTemplatesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertruletemplates.AlertRuleTemplatesClient", "AlertRuleTemplatesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertruletemplates.AlertRuleTemplatesClient", "AlertRuleTemplatesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForAlertRuleTemplatesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertruletemplates.AlertRuleTemplatesClient", "AlertRuleTemplatesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForAlertRuleTemplatesList prepares the AlertRuleTemplatesList request.
func (c AlertRuleTemplatesClient) preparerForAlertRuleTemplatesList(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.SecurityInsights/alertRuleTemplates", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForAlertRuleTemplatesListWithNextLink prepares the AlertRuleTemplatesList request with the given nextLink token.
func (c AlertRuleTemplatesClient) preparerForAlertRuleTemplatesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForAlertRuleTemplatesList handles the response to the AlertRuleTemplatesList request. The method always
// closes the http.Response Body.
func (c AlertRuleTemplatesClient) responderForAlertRuleTemplatesList(resp *http.Response) (result AlertRuleTemplatesListOperationResponse, err error) {
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
	temp := make([]AlertRuleTemplate, 0)
	for i, v := range respObj.Values {
		val, err := unmarshalAlertRuleTemplateImplementation(v)
		if err != nil {
			err = fmt.Errorf("unmarshalling item %d for AlertRuleTemplate (%q): %+v", i, v, err)
			return result, err
		}
		temp = append(temp, val)
	}
	result.Model = &temp
	result.nextLink = respObj.NextLink
	if respObj.NextLink != nil {
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result AlertRuleTemplatesListOperationResponse, err error) {
			req, err := c.preparerForAlertRuleTemplatesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "alertruletemplates.AlertRuleTemplatesClient", "AlertRuleTemplatesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "alertruletemplates.AlertRuleTemplatesClient", "AlertRuleTemplatesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForAlertRuleTemplatesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "alertruletemplates.AlertRuleTemplatesClient", "AlertRuleTemplatesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// AlertRuleTemplatesListComplete retrieves all of the results into a single object
func (c AlertRuleTemplatesClient) AlertRuleTemplatesListComplete(ctx context.Context, id WorkspaceId) (AlertRuleTemplatesListCompleteResult, error) {
	return c.AlertRuleTemplatesListCompleteMatchingPredicate(ctx, id, AlertRuleTemplateOperationPredicate{})
}

// AlertRuleTemplatesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AlertRuleTemplatesClient) AlertRuleTemplatesListCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, predicate AlertRuleTemplateOperationPredicate) (resp AlertRuleTemplatesListCompleteResult, err error) {
	items := make([]AlertRuleTemplate, 0)

	page, err := c.AlertRuleTemplatesList(ctx, id)
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

	out := AlertRuleTemplatesListCompleteResult{
		Items: items,
	}
	return out, nil
}
