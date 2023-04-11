package workflowrunactions

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

type WorkflowRunActionRequestHistoriesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]RequestHistory

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (WorkflowRunActionRequestHistoriesListOperationResponse, error)
}

type WorkflowRunActionRequestHistoriesListCompleteResult struct {
	Items []RequestHistory
}

func (r WorkflowRunActionRequestHistoriesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r WorkflowRunActionRequestHistoriesListOperationResponse) LoadMore(ctx context.Context) (resp WorkflowRunActionRequestHistoriesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// WorkflowRunActionRequestHistoriesList ...
func (c WorkflowRunActionsClient) WorkflowRunActionRequestHistoriesList(ctx context.Context, id ActionId) (resp WorkflowRunActionRequestHistoriesListOperationResponse, err error) {
	req, err := c.preparerForWorkflowRunActionRequestHistoriesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRequestHistoriesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRequestHistoriesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForWorkflowRunActionRequestHistoriesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRequestHistoriesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForWorkflowRunActionRequestHistoriesList prepares the WorkflowRunActionRequestHistoriesList request.
func (c WorkflowRunActionsClient) preparerForWorkflowRunActionRequestHistoriesList(ctx context.Context, id ActionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/requestHistories", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForWorkflowRunActionRequestHistoriesListWithNextLink prepares the WorkflowRunActionRequestHistoriesList request with the given nextLink token.
func (c WorkflowRunActionsClient) preparerForWorkflowRunActionRequestHistoriesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForWorkflowRunActionRequestHistoriesList handles the response to the WorkflowRunActionRequestHistoriesList request. The method always
// closes the http.Response Body.
func (c WorkflowRunActionsClient) responderForWorkflowRunActionRequestHistoriesList(resp *http.Response) (result WorkflowRunActionRequestHistoriesListOperationResponse, err error) {
	type page struct {
		Values   []RequestHistory `json:"value"`
		NextLink *string          `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result WorkflowRunActionRequestHistoriesListOperationResponse, err error) {
			req, err := c.preparerForWorkflowRunActionRequestHistoriesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRequestHistoriesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRequestHistoriesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForWorkflowRunActionRequestHistoriesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRequestHistoriesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// WorkflowRunActionRequestHistoriesListComplete retrieves all of the results into a single object
func (c WorkflowRunActionsClient) WorkflowRunActionRequestHistoriesListComplete(ctx context.Context, id ActionId) (WorkflowRunActionRequestHistoriesListCompleteResult, error) {
	return c.WorkflowRunActionRequestHistoriesListCompleteMatchingPredicate(ctx, id, RequestHistoryOperationPredicate{})
}

// WorkflowRunActionRequestHistoriesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c WorkflowRunActionsClient) WorkflowRunActionRequestHistoriesListCompleteMatchingPredicate(ctx context.Context, id ActionId, predicate RequestHistoryOperationPredicate) (resp WorkflowRunActionRequestHistoriesListCompleteResult, err error) {
	items := make([]RequestHistory, 0)

	page, err := c.WorkflowRunActionRequestHistoriesList(ctx, id)
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

	out := WorkflowRunActionRequestHistoriesListCompleteResult{
		Items: items,
	}
	return out, nil
}
