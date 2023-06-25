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

type WorkflowRunActionRepetitionsRequestHistoriesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]RequestHistory

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (WorkflowRunActionRepetitionsRequestHistoriesListOperationResponse, error)
}

type WorkflowRunActionRepetitionsRequestHistoriesListCompleteResult struct {
	Items []RequestHistory
}

func (r WorkflowRunActionRepetitionsRequestHistoriesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r WorkflowRunActionRepetitionsRequestHistoriesListOperationResponse) LoadMore(ctx context.Context) (resp WorkflowRunActionRepetitionsRequestHistoriesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// WorkflowRunActionRepetitionsRequestHistoriesList ...
func (c WorkflowRunActionsClient) WorkflowRunActionRepetitionsRequestHistoriesList(ctx context.Context, id RepetitionId) (resp WorkflowRunActionRepetitionsRequestHistoriesListOperationResponse, err error) {
	req, err := c.preparerForWorkflowRunActionRepetitionsRequestHistoriesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsRequestHistoriesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsRequestHistoriesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForWorkflowRunActionRepetitionsRequestHistoriesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsRequestHistoriesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForWorkflowRunActionRepetitionsRequestHistoriesList prepares the WorkflowRunActionRepetitionsRequestHistoriesList request.
func (c WorkflowRunActionsClient) preparerForWorkflowRunActionRepetitionsRequestHistoriesList(ctx context.Context, id RepetitionId) (*http.Request, error) {
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

// preparerForWorkflowRunActionRepetitionsRequestHistoriesListWithNextLink prepares the WorkflowRunActionRepetitionsRequestHistoriesList request with the given nextLink token.
func (c WorkflowRunActionsClient) preparerForWorkflowRunActionRepetitionsRequestHistoriesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForWorkflowRunActionRepetitionsRequestHistoriesList handles the response to the WorkflowRunActionRepetitionsRequestHistoriesList request. The method always
// closes the http.Response Body.
func (c WorkflowRunActionsClient) responderForWorkflowRunActionRepetitionsRequestHistoriesList(resp *http.Response) (result WorkflowRunActionRepetitionsRequestHistoriesListOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result WorkflowRunActionRepetitionsRequestHistoriesListOperationResponse, err error) {
			req, err := c.preparerForWorkflowRunActionRepetitionsRequestHistoriesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsRequestHistoriesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsRequestHistoriesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForWorkflowRunActionRepetitionsRequestHistoriesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "workflowrunactions.WorkflowRunActionsClient", "WorkflowRunActionRepetitionsRequestHistoriesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// WorkflowRunActionRepetitionsRequestHistoriesListComplete retrieves all of the results into a single object
func (c WorkflowRunActionsClient) WorkflowRunActionRepetitionsRequestHistoriesListComplete(ctx context.Context, id RepetitionId) (WorkflowRunActionRepetitionsRequestHistoriesListCompleteResult, error) {
	return c.WorkflowRunActionRepetitionsRequestHistoriesListCompleteMatchingPredicate(ctx, id, RequestHistoryOperationPredicate{})
}

// WorkflowRunActionRepetitionsRequestHistoriesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c WorkflowRunActionsClient) WorkflowRunActionRepetitionsRequestHistoriesListCompleteMatchingPredicate(ctx context.Context, id RepetitionId, predicate RequestHistoryOperationPredicate) (resp WorkflowRunActionRepetitionsRequestHistoriesListCompleteResult, err error) {
	items := make([]RequestHistory, 0)

	page, err := c.WorkflowRunActionRepetitionsRequestHistoriesList(ctx, id)
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

	out := WorkflowRunActionRepetitionsRequestHistoriesListCompleteResult{
		Items: items,
	}
	return out, nil
}
