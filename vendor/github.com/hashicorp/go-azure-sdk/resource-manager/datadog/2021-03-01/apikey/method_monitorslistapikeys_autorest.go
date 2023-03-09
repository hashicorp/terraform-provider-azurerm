package apikey

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

type MonitorsListApiKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DatadogApiKey

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (MonitorsListApiKeysOperationResponse, error)
}

type MonitorsListApiKeysCompleteResult struct {
	Items []DatadogApiKey
}

func (r MonitorsListApiKeysOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r MonitorsListApiKeysOperationResponse) LoadMore(ctx context.Context) (resp MonitorsListApiKeysOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// MonitorsListApiKeys ...
func (c ApiKeyClient) MonitorsListApiKeys(ctx context.Context, id MonitorId) (resp MonitorsListApiKeysOperationResponse, err error) {
	req, err := c.preparerForMonitorsListApiKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apikey.ApiKeyClient", "MonitorsListApiKeys", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "apikey.ApiKeyClient", "MonitorsListApiKeys", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForMonitorsListApiKeys(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apikey.ApiKeyClient", "MonitorsListApiKeys", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForMonitorsListApiKeys prepares the MonitorsListApiKeys request.
func (c ApiKeyClient) preparerForMonitorsListApiKeys(ctx context.Context, id MonitorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listApiKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForMonitorsListApiKeysWithNextLink prepares the MonitorsListApiKeys request with the given nextLink token.
func (c ApiKeyClient) preparerForMonitorsListApiKeysWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMonitorsListApiKeys handles the response to the MonitorsListApiKeys request. The method always
// closes the http.Response Body.
func (c ApiKeyClient) responderForMonitorsListApiKeys(resp *http.Response) (result MonitorsListApiKeysOperationResponse, err error) {
	type page struct {
		Values   []DatadogApiKey `json:"value"`
		NextLink *string         `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result MonitorsListApiKeysOperationResponse, err error) {
			req, err := c.preparerForMonitorsListApiKeysWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "apikey.ApiKeyClient", "MonitorsListApiKeys", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "apikey.ApiKeyClient", "MonitorsListApiKeys", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForMonitorsListApiKeys(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "apikey.ApiKeyClient", "MonitorsListApiKeys", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// MonitorsListApiKeysComplete retrieves all of the results into a single object
func (c ApiKeyClient) MonitorsListApiKeysComplete(ctx context.Context, id MonitorId) (MonitorsListApiKeysCompleteResult, error) {
	return c.MonitorsListApiKeysCompleteMatchingPredicate(ctx, id, DatadogApiKeyOperationPredicate{})
}

// MonitorsListApiKeysCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ApiKeyClient) MonitorsListApiKeysCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate DatadogApiKeyOperationPredicate) (resp MonitorsListApiKeysCompleteResult, err error) {
	items := make([]DatadogApiKey, 0)

	page, err := c.MonitorsListApiKeys(ctx, id)
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

	out := MonitorsListApiKeysCompleteResult{
		Items: items,
	}
	return out, nil
}
