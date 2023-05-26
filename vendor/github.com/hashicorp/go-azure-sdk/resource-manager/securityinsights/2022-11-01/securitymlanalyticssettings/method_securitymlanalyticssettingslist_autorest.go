package securitymlanalyticssettings

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

type SecurityMLAnalyticsSettingsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]SecurityMLAnalyticsSetting

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (SecurityMLAnalyticsSettingsListOperationResponse, error)
}

type SecurityMLAnalyticsSettingsListCompleteResult struct {
	Items []SecurityMLAnalyticsSetting
}

func (r SecurityMLAnalyticsSettingsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r SecurityMLAnalyticsSettingsListOperationResponse) LoadMore(ctx context.Context) (resp SecurityMLAnalyticsSettingsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// SecurityMLAnalyticsSettingsList ...
func (c SecurityMLAnalyticsSettingsClient) SecurityMLAnalyticsSettingsList(ctx context.Context, id WorkspaceId) (resp SecurityMLAnalyticsSettingsListOperationResponse, err error) {
	req, err := c.preparerForSecurityMLAnalyticsSettingsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient", "SecurityMLAnalyticsSettingsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient", "SecurityMLAnalyticsSettingsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForSecurityMLAnalyticsSettingsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient", "SecurityMLAnalyticsSettingsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForSecurityMLAnalyticsSettingsList prepares the SecurityMLAnalyticsSettingsList request.
func (c SecurityMLAnalyticsSettingsClient) preparerForSecurityMLAnalyticsSettingsList(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.SecurityInsights/securityMLAnalyticsSettings", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForSecurityMLAnalyticsSettingsListWithNextLink prepares the SecurityMLAnalyticsSettingsList request with the given nextLink token.
func (c SecurityMLAnalyticsSettingsClient) preparerForSecurityMLAnalyticsSettingsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForSecurityMLAnalyticsSettingsList handles the response to the SecurityMLAnalyticsSettingsList request. The method always
// closes the http.Response Body.
func (c SecurityMLAnalyticsSettingsClient) responderForSecurityMLAnalyticsSettingsList(resp *http.Response) (result SecurityMLAnalyticsSettingsListOperationResponse, err error) {
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
	temp := make([]SecurityMLAnalyticsSetting, 0)
	for i, v := range respObj.Values {
		val, err := unmarshalSecurityMLAnalyticsSettingImplementation(v)
		if err != nil {
			err = fmt.Errorf("unmarshalling item %d for SecurityMLAnalyticsSetting (%q): %+v", i, v, err)
			return result, err
		}
		temp = append(temp, val)
	}
	result.Model = &temp
	result.nextLink = respObj.NextLink
	if respObj.NextLink != nil {
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result SecurityMLAnalyticsSettingsListOperationResponse, err error) {
			req, err := c.preparerForSecurityMLAnalyticsSettingsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient", "SecurityMLAnalyticsSettingsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient", "SecurityMLAnalyticsSettingsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForSecurityMLAnalyticsSettingsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient", "SecurityMLAnalyticsSettingsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// SecurityMLAnalyticsSettingsListComplete retrieves all of the results into a single object
func (c SecurityMLAnalyticsSettingsClient) SecurityMLAnalyticsSettingsListComplete(ctx context.Context, id WorkspaceId) (SecurityMLAnalyticsSettingsListCompleteResult, error) {
	return c.SecurityMLAnalyticsSettingsListCompleteMatchingPredicate(ctx, id, SecurityMLAnalyticsSettingOperationPredicate{})
}

// SecurityMLAnalyticsSettingsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c SecurityMLAnalyticsSettingsClient) SecurityMLAnalyticsSettingsListCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, predicate SecurityMLAnalyticsSettingOperationPredicate) (resp SecurityMLAnalyticsSettingsListCompleteResult, err error) {
	items := make([]SecurityMLAnalyticsSetting, 0)

	page, err := c.SecurityMLAnalyticsSettingsList(ctx, id)
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

	out := SecurityMLAnalyticsSettingsListCompleteResult{
		Items: items,
	}
	return out, nil
}
