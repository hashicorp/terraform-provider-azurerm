package deviceupdates

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

type InstancesListByAccountOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Instance

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (InstancesListByAccountOperationResponse, error)
}

type InstancesListByAccountCompleteResult struct {
	Items []Instance
}

func (r InstancesListByAccountOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r InstancesListByAccountOperationResponse) LoadMore(ctx context.Context) (resp InstancesListByAccountOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// InstancesListByAccount ...
func (c DeviceupdatesClient) InstancesListByAccount(ctx context.Context, id AccountId) (resp InstancesListByAccountOperationResponse, err error) {
	req, err := c.preparerForInstancesListByAccount(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesListByAccount", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesListByAccount", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForInstancesListByAccount(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesListByAccount", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForInstancesListByAccount prepares the InstancesListByAccount request.
func (c DeviceupdatesClient) preparerForInstancesListByAccount(ctx context.Context, id AccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/instances", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForInstancesListByAccountWithNextLink prepares the InstancesListByAccount request with the given nextLink token.
func (c DeviceupdatesClient) preparerForInstancesListByAccountWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForInstancesListByAccount handles the response to the InstancesListByAccount request. The method always
// closes the http.Response Body.
func (c DeviceupdatesClient) responderForInstancesListByAccount(resp *http.Response) (result InstancesListByAccountOperationResponse, err error) {
	type page struct {
		Values   []Instance `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result InstancesListByAccountOperationResponse, err error) {
			req, err := c.preparerForInstancesListByAccountWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesListByAccount", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesListByAccount", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForInstancesListByAccount(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "deviceupdates.DeviceupdatesClient", "InstancesListByAccount", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// InstancesListByAccountComplete retrieves all of the results into a single object
func (c DeviceupdatesClient) InstancesListByAccountComplete(ctx context.Context, id AccountId) (InstancesListByAccountCompleteResult, error) {
	return c.InstancesListByAccountCompleteMatchingPredicate(ctx, id, InstanceOperationPredicate{})
}

// InstancesListByAccountCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c DeviceupdatesClient) InstancesListByAccountCompleteMatchingPredicate(ctx context.Context, id AccountId, predicate InstanceOperationPredicate) (resp InstancesListByAccountCompleteResult, err error) {
	items := make([]Instance, 0)

	page, err := c.InstancesListByAccount(ctx, id)
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

	out := InstancesListByAccountCompleteResult{
		Items: items,
	}
	return out, nil
}
