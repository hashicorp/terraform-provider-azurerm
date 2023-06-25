package updatesummaries

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

type UpdateSummariesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]UpdateSummaries

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (UpdateSummariesListOperationResponse, error)
}

type UpdateSummariesListCompleteResult struct {
	Items []UpdateSummaries
}

func (r UpdateSummariesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r UpdateSummariesListOperationResponse) LoadMore(ctx context.Context) (resp UpdateSummariesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// UpdateSummariesList ...
func (c UpdateSummariesClient) UpdateSummariesList(ctx context.Context, id ClusterId) (resp UpdateSummariesListOperationResponse, err error) {
	req, err := c.preparerForUpdateSummariesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updatesummaries.UpdateSummariesClient", "UpdateSummariesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "updatesummaries.UpdateSummariesClient", "UpdateSummariesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForUpdateSummariesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updatesummaries.UpdateSummariesClient", "UpdateSummariesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForUpdateSummariesList prepares the UpdateSummariesList request.
func (c UpdateSummariesClient) preparerForUpdateSummariesList(ctx context.Context, id ClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/updateSummaries", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForUpdateSummariesListWithNextLink prepares the UpdateSummariesList request with the given nextLink token.
func (c UpdateSummariesClient) preparerForUpdateSummariesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForUpdateSummariesList handles the response to the UpdateSummariesList request. The method always
// closes the http.Response Body.
func (c UpdateSummariesClient) responderForUpdateSummariesList(resp *http.Response) (result UpdateSummariesListOperationResponse, err error) {
	type page struct {
		Values   []UpdateSummaries `json:"value"`
		NextLink *string           `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result UpdateSummariesListOperationResponse, err error) {
			req, err := c.preparerForUpdateSummariesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "updatesummaries.UpdateSummariesClient", "UpdateSummariesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "updatesummaries.UpdateSummariesClient", "UpdateSummariesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForUpdateSummariesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "updatesummaries.UpdateSummariesClient", "UpdateSummariesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// UpdateSummariesListComplete retrieves all of the results into a single object
func (c UpdateSummariesClient) UpdateSummariesListComplete(ctx context.Context, id ClusterId) (UpdateSummariesListCompleteResult, error) {
	return c.UpdateSummariesListCompleteMatchingPredicate(ctx, id, UpdateSummariesOperationPredicate{})
}

// UpdateSummariesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c UpdateSummariesClient) UpdateSummariesListCompleteMatchingPredicate(ctx context.Context, id ClusterId, predicate UpdateSummariesOperationPredicate) (resp UpdateSummariesListCompleteResult, err error) {
	items := make([]UpdateSummaries, 0)

	page, err := c.UpdateSummariesList(ctx, id)
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

	out := UpdateSummariesListCompleteResult{
		Items: items,
	}
	return out, nil
}
