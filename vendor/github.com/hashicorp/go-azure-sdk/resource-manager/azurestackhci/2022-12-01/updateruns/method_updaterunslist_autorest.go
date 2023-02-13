package updateruns

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

type UpdateRunsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]UpdateRun

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (UpdateRunsListOperationResponse, error)
}

type UpdateRunsListCompleteResult struct {
	Items []UpdateRun
}

func (r UpdateRunsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r UpdateRunsListOperationResponse) LoadMore(ctx context.Context) (resp UpdateRunsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// UpdateRunsList ...
func (c UpdateRunsClient) UpdateRunsList(ctx context.Context, id UpdateId) (resp UpdateRunsListOperationResponse, err error) {
	req, err := c.preparerForUpdateRunsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updateruns.UpdateRunsClient", "UpdateRunsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "updateruns.UpdateRunsClient", "UpdateRunsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForUpdateRunsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updateruns.UpdateRunsClient", "UpdateRunsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForUpdateRunsList prepares the UpdateRunsList request.
func (c UpdateRunsClient) preparerForUpdateRunsList(ctx context.Context, id UpdateId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/updateRuns", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForUpdateRunsListWithNextLink prepares the UpdateRunsList request with the given nextLink token.
func (c UpdateRunsClient) preparerForUpdateRunsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForUpdateRunsList handles the response to the UpdateRunsList request. The method always
// closes the http.Response Body.
func (c UpdateRunsClient) responderForUpdateRunsList(resp *http.Response) (result UpdateRunsListOperationResponse, err error) {
	type page struct {
		Values   []UpdateRun `json:"value"`
		NextLink *string     `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result UpdateRunsListOperationResponse, err error) {
			req, err := c.preparerForUpdateRunsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "updateruns.UpdateRunsClient", "UpdateRunsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "updateruns.UpdateRunsClient", "UpdateRunsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForUpdateRunsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "updateruns.UpdateRunsClient", "UpdateRunsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// UpdateRunsListComplete retrieves all of the results into a single object
func (c UpdateRunsClient) UpdateRunsListComplete(ctx context.Context, id UpdateId) (UpdateRunsListCompleteResult, error) {
	return c.UpdateRunsListCompleteMatchingPredicate(ctx, id, UpdateRunOperationPredicate{})
}

// UpdateRunsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c UpdateRunsClient) UpdateRunsListCompleteMatchingPredicate(ctx context.Context, id UpdateId, predicate UpdateRunOperationPredicate) (resp UpdateRunsListCompleteResult, err error) {
	items := make([]UpdateRun, 0)

	page, err := c.UpdateRunsList(ctx, id)
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

	out := UpdateRunsListCompleteResult{
		Items: items,
	}
	return out, nil
}
