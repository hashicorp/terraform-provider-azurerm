package alertsmanagement

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

type AlertProcessingRulesListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AlertProcessingRule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (AlertProcessingRulesListByResourceGroupOperationResponse, error)
}

type AlertProcessingRulesListByResourceGroupCompleteResult struct {
	Items []AlertProcessingRule
}

func (r AlertProcessingRulesListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r AlertProcessingRulesListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp AlertProcessingRulesListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// AlertProcessingRulesListByResourceGroup ...
func (c AlertsManagementClient) AlertProcessingRulesListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp AlertProcessingRulesListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForAlertProcessingRulesListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.AlertsManagementClient", "AlertProcessingRulesListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.AlertsManagementClient", "AlertProcessingRulesListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForAlertProcessingRulesListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.AlertsManagementClient", "AlertProcessingRulesListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// AlertProcessingRulesListByResourceGroupComplete retrieves all of the results into a single object
func (c AlertsManagementClient) AlertProcessingRulesListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (AlertProcessingRulesListByResourceGroupCompleteResult, error) {
	return c.AlertProcessingRulesListByResourceGroupCompleteMatchingPredicate(ctx, id, AlertProcessingRuleOperationPredicate{})
}

// AlertProcessingRulesListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AlertsManagementClient) AlertProcessingRulesListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate AlertProcessingRuleOperationPredicate) (resp AlertProcessingRulesListByResourceGroupCompleteResult, err error) {
	items := make([]AlertProcessingRule, 0)

	page, err := c.AlertProcessingRulesListByResourceGroup(ctx, id)
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

	out := AlertProcessingRulesListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForAlertProcessingRulesListByResourceGroup prepares the AlertProcessingRulesListByResourceGroup request.
func (c AlertsManagementClient) preparerForAlertProcessingRulesListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.AlertsManagement/actionRules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForAlertProcessingRulesListByResourceGroupWithNextLink prepares the AlertProcessingRulesListByResourceGroup request with the given nextLink token.
func (c AlertsManagementClient) preparerForAlertProcessingRulesListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForAlertProcessingRulesListByResourceGroup handles the response to the AlertProcessingRulesListByResourceGroup request. The method always
// closes the http.Response Body.
func (c AlertsManagementClient) responderForAlertProcessingRulesListByResourceGroup(resp *http.Response) (result AlertProcessingRulesListByResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []AlertProcessingRule `json:"value"`
		NextLink *string               `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result AlertProcessingRulesListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForAlertProcessingRulesListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "alertsmanagement.AlertsManagementClient", "AlertProcessingRulesListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "alertsmanagement.AlertsManagementClient", "AlertProcessingRulesListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForAlertProcessingRulesListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "alertsmanagement.AlertsManagementClient", "AlertProcessingRulesListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
