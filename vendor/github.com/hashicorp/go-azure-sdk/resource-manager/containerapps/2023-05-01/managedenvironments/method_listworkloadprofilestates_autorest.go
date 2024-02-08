package managedenvironments

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

type ListWorkloadProfileStatesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]WorkloadProfileStates

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListWorkloadProfileStatesOperationResponse, error)
}

type ListWorkloadProfileStatesCompleteResult struct {
	Items []WorkloadProfileStates
}

func (r ListWorkloadProfileStatesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListWorkloadProfileStatesOperationResponse) LoadMore(ctx context.Context) (resp ListWorkloadProfileStatesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListWorkloadProfileStates ...
func (c ManagedEnvironmentsClient) ListWorkloadProfileStates(ctx context.Context, id ManagedEnvironmentId) (resp ListWorkloadProfileStatesOperationResponse, err error) {
	req, err := c.preparerForListWorkloadProfileStates(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ListWorkloadProfileStates", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ListWorkloadProfileStates", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListWorkloadProfileStates(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ListWorkloadProfileStates", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListWorkloadProfileStates prepares the ListWorkloadProfileStates request.
func (c ManagedEnvironmentsClient) preparerForListWorkloadProfileStates(ctx context.Context, id ManagedEnvironmentId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/workloadProfileStates", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListWorkloadProfileStatesWithNextLink prepares the ListWorkloadProfileStates request with the given nextLink token.
func (c ManagedEnvironmentsClient) preparerForListWorkloadProfileStatesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListWorkloadProfileStates handles the response to the ListWorkloadProfileStates request. The method always
// closes the http.Response Body.
func (c ManagedEnvironmentsClient) responderForListWorkloadProfileStates(resp *http.Response) (result ListWorkloadProfileStatesOperationResponse, err error) {
	type page struct {
		Values   []WorkloadProfileStates `json:"value"`
		NextLink *string                 `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListWorkloadProfileStatesOperationResponse, err error) {
			req, err := c.preparerForListWorkloadProfileStatesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ListWorkloadProfileStates", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ListWorkloadProfileStates", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListWorkloadProfileStates(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedenvironments.ManagedEnvironmentsClient", "ListWorkloadProfileStates", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListWorkloadProfileStatesComplete retrieves all of the results into a single object
func (c ManagedEnvironmentsClient) ListWorkloadProfileStatesComplete(ctx context.Context, id ManagedEnvironmentId) (ListWorkloadProfileStatesCompleteResult, error) {
	return c.ListWorkloadProfileStatesCompleteMatchingPredicate(ctx, id, WorkloadProfileStatesOperationPredicate{})
}

// ListWorkloadProfileStatesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ManagedEnvironmentsClient) ListWorkloadProfileStatesCompleteMatchingPredicate(ctx context.Context, id ManagedEnvironmentId, predicate WorkloadProfileStatesOperationPredicate) (resp ListWorkloadProfileStatesCompleteResult, err error) {
	items := make([]WorkloadProfileStates, 0)

	page, err := c.ListWorkloadProfileStates(ctx, id)
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

	out := ListWorkloadProfileStatesCompleteResult{
		Items: items,
	}
	return out, nil
}
