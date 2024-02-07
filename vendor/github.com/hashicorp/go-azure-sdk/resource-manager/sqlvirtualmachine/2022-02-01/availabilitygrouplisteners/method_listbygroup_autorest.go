package availabilitygrouplisteners

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

type ListByGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AvailabilityGroupListener

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByGroupOperationResponse, error)
}

type ListByGroupCompleteResult struct {
	Items []AvailabilityGroupListener
}

func (r ListByGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByGroupOperationResponse) LoadMore(ctx context.Context) (resp ListByGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByGroup ...
func (c AvailabilityGroupListenersClient) ListByGroup(ctx context.Context, id SqlVirtualMachineGroupId) (resp ListByGroupOperationResponse, err error) {
	req, err := c.preparerForListByGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "availabilitygrouplisteners.AvailabilityGroupListenersClient", "ListByGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "availabilitygrouplisteners.AvailabilityGroupListenersClient", "ListByGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "availabilitygrouplisteners.AvailabilityGroupListenersClient", "ListByGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByGroup prepares the ListByGroup request.
func (c AvailabilityGroupListenersClient) preparerForListByGroup(ctx context.Context, id SqlVirtualMachineGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/availabilityGroupListeners", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByGroupWithNextLink prepares the ListByGroup request with the given nextLink token.
func (c AvailabilityGroupListenersClient) preparerForListByGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByGroup handles the response to the ListByGroup request. The method always
// closes the http.Response Body.
func (c AvailabilityGroupListenersClient) responderForListByGroup(resp *http.Response) (result ListByGroupOperationResponse, err error) {
	type page struct {
		Values   []AvailabilityGroupListener `json:"value"`
		NextLink *string                     `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByGroupOperationResponse, err error) {
			req, err := c.preparerForListByGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "availabilitygrouplisteners.AvailabilityGroupListenersClient", "ListByGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "availabilitygrouplisteners.AvailabilityGroupListenersClient", "ListByGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "availabilitygrouplisteners.AvailabilityGroupListenersClient", "ListByGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByGroupComplete retrieves all of the results into a single object
func (c AvailabilityGroupListenersClient) ListByGroupComplete(ctx context.Context, id SqlVirtualMachineGroupId) (ListByGroupCompleteResult, error) {
	return c.ListByGroupCompleteMatchingPredicate(ctx, id, AvailabilityGroupListenerOperationPredicate{})
}

// ListByGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AvailabilityGroupListenersClient) ListByGroupCompleteMatchingPredicate(ctx context.Context, id SqlVirtualMachineGroupId, predicate AvailabilityGroupListenerOperationPredicate) (resp ListByGroupCompleteResult, err error) {
	items := make([]AvailabilityGroupListener, 0)

	page, err := c.ListByGroup(ctx, id)
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

	out := ListByGroupCompleteResult{
		Items: items,
	}
	return out, nil
}
