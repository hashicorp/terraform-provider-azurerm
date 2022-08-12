package applicationinsights

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

type WorkbooksListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Workbook

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (WorkbooksListByResourceGroupOperationResponse, error)
}

type WorkbooksListByResourceGroupCompleteResult struct {
	Items []Workbook
}

func (r WorkbooksListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r WorkbooksListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp WorkbooksListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type WorkbooksListByResourceGroupOperationOptions struct {
	CanFetchContent *bool
	Category        *CategoryType
	SourceId        *string
	Tags            *string
}

func DefaultWorkbooksListByResourceGroupOperationOptions() WorkbooksListByResourceGroupOperationOptions {
	return WorkbooksListByResourceGroupOperationOptions{}
}

func (o WorkbooksListByResourceGroupOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o WorkbooksListByResourceGroupOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.CanFetchContent != nil {
		out["canFetchContent"] = *o.CanFetchContent
	}

	if o.Category != nil {
		out["category"] = *o.Category
	}

	if o.SourceId != nil {
		out["sourceId"] = *o.SourceId
	}

	if o.Tags != nil {
		out["tags"] = *o.Tags
	}

	return out
}

// WorkbooksListByResourceGroup ...
func (c ApplicationInsightsClient) WorkbooksListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options WorkbooksListByResourceGroupOperationOptions) (resp WorkbooksListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForWorkbooksListByResourceGroup(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForWorkbooksListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForWorkbooksListByResourceGroup prepares the WorkbooksListByResourceGroup request.
func (c ApplicationInsightsClient) preparerForWorkbooksListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options WorkbooksListByResourceGroupOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Insights/workbooks", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForWorkbooksListByResourceGroupWithNextLink prepares the WorkbooksListByResourceGroup request with the given nextLink token.
func (c ApplicationInsightsClient) preparerForWorkbooksListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForWorkbooksListByResourceGroup handles the response to the WorkbooksListByResourceGroup request. The method always
// closes the http.Response Body.
func (c ApplicationInsightsClient) responderForWorkbooksListByResourceGroup(resp *http.Response) (result WorkbooksListByResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []Workbook `json:"value"`
		NextLink *string    `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result WorkbooksListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForWorkbooksListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForWorkbooksListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// WorkbooksListByResourceGroupComplete retrieves all of the results into a single object
func (c ApplicationInsightsClient) WorkbooksListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId, options WorkbooksListByResourceGroupOperationOptions) (WorkbooksListByResourceGroupCompleteResult, error) {
	return c.WorkbooksListByResourceGroupCompleteMatchingPredicate(ctx, id, options, WorkbookOperationPredicate{})
}

// WorkbooksListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ApplicationInsightsClient) WorkbooksListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, options WorkbooksListByResourceGroupOperationOptions, predicate WorkbookOperationPredicate) (resp WorkbooksListByResourceGroupCompleteResult, err error) {
	items := make([]Workbook, 0)

	page, err := c.WorkbooksListByResourceGroup(ctx, id, options)
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

	out := WorkbooksListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}
