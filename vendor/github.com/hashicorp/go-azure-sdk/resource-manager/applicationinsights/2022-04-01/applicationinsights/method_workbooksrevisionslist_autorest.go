package applicationinsights

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

type WorkbooksRevisionsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Workbook

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (WorkbooksRevisionsListOperationResponse, error)
}

type WorkbooksRevisionsListCompleteResult struct {
	Items []Workbook
}

func (r WorkbooksRevisionsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r WorkbooksRevisionsListOperationResponse) LoadMore(ctx context.Context) (resp WorkbooksRevisionsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// WorkbooksRevisionsList ...
func (c ApplicationInsightsClient) WorkbooksRevisionsList(ctx context.Context, id WorkbookId) (resp WorkbooksRevisionsListOperationResponse, err error) {
	req, err := c.preparerForWorkbooksRevisionsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksRevisionsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksRevisionsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForWorkbooksRevisionsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksRevisionsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForWorkbooksRevisionsList prepares the WorkbooksRevisionsList request.
func (c ApplicationInsightsClient) preparerForWorkbooksRevisionsList(ctx context.Context, id WorkbookId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/revisions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForWorkbooksRevisionsListWithNextLink prepares the WorkbooksRevisionsList request with the given nextLink token.
func (c ApplicationInsightsClient) preparerForWorkbooksRevisionsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForWorkbooksRevisionsList handles the response to the WorkbooksRevisionsList request. The method always
// closes the http.Response Body.
func (c ApplicationInsightsClient) responderForWorkbooksRevisionsList(resp *http.Response) (result WorkbooksRevisionsListOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result WorkbooksRevisionsListOperationResponse, err error) {
			req, err := c.preparerForWorkbooksRevisionsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksRevisionsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksRevisionsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForWorkbooksRevisionsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksRevisionsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// WorkbooksRevisionsListComplete retrieves all of the results into a single object
func (c ApplicationInsightsClient) WorkbooksRevisionsListComplete(ctx context.Context, id WorkbookId) (WorkbooksRevisionsListCompleteResult, error) {
	return c.WorkbooksRevisionsListCompleteMatchingPredicate(ctx, id, WorkbookOperationPredicate{})
}

// WorkbooksRevisionsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ApplicationInsightsClient) WorkbooksRevisionsListCompleteMatchingPredicate(ctx context.Context, id WorkbookId, predicate WorkbookOperationPredicate) (resp WorkbooksRevisionsListCompleteResult, err error) {
	items := make([]Workbook, 0)

	page, err := c.WorkbooksRevisionsList(ctx, id)
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

	out := WorkbooksRevisionsListCompleteResult{
		Items: items,
	}
	return out, nil
}
