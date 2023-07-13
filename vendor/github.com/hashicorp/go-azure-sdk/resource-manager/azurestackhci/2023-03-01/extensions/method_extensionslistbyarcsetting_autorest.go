package extensions

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

type ExtensionsListByArcSettingOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Extension

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ExtensionsListByArcSettingOperationResponse, error)
}

type ExtensionsListByArcSettingCompleteResult struct {
	Items []Extension
}

func (r ExtensionsListByArcSettingOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ExtensionsListByArcSettingOperationResponse) LoadMore(ctx context.Context) (resp ExtensionsListByArcSettingOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ExtensionsListByArcSetting ...
func (c ExtensionsClient) ExtensionsListByArcSetting(ctx context.Context, id ArcSettingId) (resp ExtensionsListByArcSettingOperationResponse, err error) {
	req, err := c.preparerForExtensionsListByArcSetting(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "extensions.ExtensionsClient", "ExtensionsListByArcSetting", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "extensions.ExtensionsClient", "ExtensionsListByArcSetting", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForExtensionsListByArcSetting(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "extensions.ExtensionsClient", "ExtensionsListByArcSetting", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForExtensionsListByArcSetting prepares the ExtensionsListByArcSetting request.
func (c ExtensionsClient) preparerForExtensionsListByArcSetting(ctx context.Context, id ArcSettingId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/extensions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForExtensionsListByArcSettingWithNextLink prepares the ExtensionsListByArcSetting request with the given nextLink token.
func (c ExtensionsClient) preparerForExtensionsListByArcSettingWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForExtensionsListByArcSetting handles the response to the ExtensionsListByArcSetting request. The method always
// closes the http.Response Body.
func (c ExtensionsClient) responderForExtensionsListByArcSetting(resp *http.Response) (result ExtensionsListByArcSettingOperationResponse, err error) {
	type page struct {
		Values   []Extension `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ExtensionsListByArcSettingOperationResponse, err error) {
			req, err := c.preparerForExtensionsListByArcSettingWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "extensions.ExtensionsClient", "ExtensionsListByArcSetting", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "extensions.ExtensionsClient", "ExtensionsListByArcSetting", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForExtensionsListByArcSetting(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "extensions.ExtensionsClient", "ExtensionsListByArcSetting", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ExtensionsListByArcSettingComplete retrieves all of the results into a single object
func (c ExtensionsClient) ExtensionsListByArcSettingComplete(ctx context.Context, id ArcSettingId) (ExtensionsListByArcSettingCompleteResult, error) {
	return c.ExtensionsListByArcSettingCompleteMatchingPredicate(ctx, id, ExtensionOperationPredicate{})
}

// ExtensionsListByArcSettingCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ExtensionsClient) ExtensionsListByArcSettingCompleteMatchingPredicate(ctx context.Context, id ArcSettingId, predicate ExtensionOperationPredicate) (resp ExtensionsListByArcSettingCompleteResult, err error) {
	items := make([]Extension, 0)

	page, err := c.ExtensionsListByArcSetting(ctx, id)
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

	out := ExtensionsListByArcSettingCompleteResult{
		Items: items,
	}
	return out, nil
}
