package synchronizationsetting

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

type ListByShareOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]SynchronizationSetting

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByShareOperationResponse, error)
}

type ListByShareCompleteResult struct {
	Items []SynchronizationSetting
}

func (r ListByShareOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByShareOperationResponse) LoadMore(ctx context.Context) (resp ListByShareOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByShare ...
func (c SynchronizationSettingClient) ListByShare(ctx context.Context, id ShareId) (resp ListByShareOperationResponse, err error) {
	req, err := c.preparerForListByShare(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "synchronizationsetting.SynchronizationSettingClient", "ListByShare", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "synchronizationsetting.SynchronizationSettingClient", "ListByShare", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByShare(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "synchronizationsetting.SynchronizationSettingClient", "ListByShare", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByShare prepares the ListByShare request.
func (c SynchronizationSettingClient) preparerForListByShare(ctx context.Context, id ShareId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/synchronizationSettings", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByShareWithNextLink prepares the ListByShare request with the given nextLink token.
func (c SynchronizationSettingClient) preparerForListByShareWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByShare handles the response to the ListByShare request. The method always
// closes the http.Response Body.
func (c SynchronizationSettingClient) responderForListByShare(resp *http.Response) (result ListByShareOperationResponse, err error) {
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
	temp := make([]SynchronizationSetting, 0)
	for i, v := range respObj.Values {
		val, err := unmarshalSynchronizationSettingImplementation(v)
		if err != nil {
			err = fmt.Errorf("unmarshalling item %d for SynchronizationSetting (%q): %+v", i, v, err)
			return result, err
		}
		temp = append(temp, val)
	}
	result.Model = &temp
	result.nextLink = respObj.NextLink
	if respObj.NextLink != nil {
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByShareOperationResponse, err error) {
			req, err := c.preparerForListByShareWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "synchronizationsetting.SynchronizationSettingClient", "ListByShare", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "synchronizationsetting.SynchronizationSettingClient", "ListByShare", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByShare(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "synchronizationsetting.SynchronizationSettingClient", "ListByShare", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByShareComplete retrieves all of the results into a single object
func (c SynchronizationSettingClient) ListByShareComplete(ctx context.Context, id ShareId) (ListByShareCompleteResult, error) {
	return c.ListByShareCompleteMatchingPredicate(ctx, id, SynchronizationSettingOperationPredicate{})
}

// ListByShareCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c SynchronizationSettingClient) ListByShareCompleteMatchingPredicate(ctx context.Context, id ShareId, predicate SynchronizationSettingOperationPredicate) (resp ListByShareCompleteResult, err error) {
	items := make([]SynchronizationSetting, 0)

	page, err := c.ListByShare(ctx, id)
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

	out := ListByShareCompleteResult{
		Items: items,
	}
	return out, nil
}
