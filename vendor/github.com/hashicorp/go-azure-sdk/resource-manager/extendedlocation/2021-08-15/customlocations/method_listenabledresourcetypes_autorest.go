package customlocations

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

type ListEnabledResourceTypesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]EnabledResourceType

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListEnabledResourceTypesOperationResponse, error)
}

type ListEnabledResourceTypesCompleteResult struct {
	Items []EnabledResourceType
}

func (r ListEnabledResourceTypesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListEnabledResourceTypesOperationResponse) LoadMore(ctx context.Context) (resp ListEnabledResourceTypesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListEnabledResourceTypes ...
func (c CustomLocationsClient) ListEnabledResourceTypes(ctx context.Context, id CustomLocationId) (resp ListEnabledResourceTypesOperationResponse, err error) {
	req, err := c.preparerForListEnabledResourceTypes(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customlocations.CustomLocationsClient", "ListEnabledResourceTypes", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "customlocations.CustomLocationsClient", "ListEnabledResourceTypes", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListEnabledResourceTypes(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "customlocations.CustomLocationsClient", "ListEnabledResourceTypes", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListEnabledResourceTypes prepares the ListEnabledResourceTypes request.
func (c CustomLocationsClient) preparerForListEnabledResourceTypes(ctx context.Context, id CustomLocationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/enabledResourceTypes", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListEnabledResourceTypesWithNextLink prepares the ListEnabledResourceTypes request with the given nextLink token.
func (c CustomLocationsClient) preparerForListEnabledResourceTypesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListEnabledResourceTypes handles the response to the ListEnabledResourceTypes request. The method always
// closes the http.Response Body.
func (c CustomLocationsClient) responderForListEnabledResourceTypes(resp *http.Response) (result ListEnabledResourceTypesOperationResponse, err error) {
	type page struct {
		Values   []EnabledResourceType `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListEnabledResourceTypesOperationResponse, err error) {
			req, err := c.preparerForListEnabledResourceTypesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "customlocations.CustomLocationsClient", "ListEnabledResourceTypes", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "customlocations.CustomLocationsClient", "ListEnabledResourceTypes", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListEnabledResourceTypes(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "customlocations.CustomLocationsClient", "ListEnabledResourceTypes", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListEnabledResourceTypesComplete retrieves all of the results into a single object
func (c CustomLocationsClient) ListEnabledResourceTypesComplete(ctx context.Context, id CustomLocationId) (ListEnabledResourceTypesCompleteResult, error) {
	return c.ListEnabledResourceTypesCompleteMatchingPredicate(ctx, id, EnabledResourceTypeOperationPredicate{})
}

// ListEnabledResourceTypesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c CustomLocationsClient) ListEnabledResourceTypesCompleteMatchingPredicate(ctx context.Context, id CustomLocationId, predicate EnabledResourceTypeOperationPredicate) (resp ListEnabledResourceTypesCompleteResult, err error) {
	items := make([]EnabledResourceType, 0)

	page, err := c.ListEnabledResourceTypes(ctx, id)
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

	out := ListEnabledResourceTypesCompleteResult{
		Items: items,
	}
	return out, nil
}
