package user

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

type ListByLabOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]User

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByLabOperationResponse, error)
}

type ListByLabCompleteResult struct {
	Items []User
}

func (r ListByLabOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByLabOperationResponse) LoadMore(ctx context.Context) (resp ListByLabOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByLab ...
func (c UserClient) ListByLab(ctx context.Context, id LabId) (resp ListByLabOperationResponse, err error) {
	req, err := c.preparerForListByLab(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "user.UserClient", "ListByLab", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "user.UserClient", "ListByLab", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByLab(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "user.UserClient", "ListByLab", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByLab prepares the ListByLab request.
func (c UserClient) preparerForListByLab(ctx context.Context, id LabId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/users", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByLabWithNextLink prepares the ListByLab request with the given nextLink token.
func (c UserClient) preparerForListByLabWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByLab handles the response to the ListByLab request. The method always
// closes the http.Response Body.
func (c UserClient) responderForListByLab(resp *http.Response) (result ListByLabOperationResponse, err error) {
	type page struct {
		Values   []User  `json:"value"`
		NextLink *string `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByLabOperationResponse, err error) {
			req, err := c.preparerForListByLabWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "user.UserClient", "ListByLab", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "user.UserClient", "ListByLab", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByLab(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "user.UserClient", "ListByLab", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByLabComplete retrieves all of the results into a single object
func (c UserClient) ListByLabComplete(ctx context.Context, id LabId) (ListByLabCompleteResult, error) {
	return c.ListByLabCompleteMatchingPredicate(ctx, id, UserOperationPredicate{})
}

// ListByLabCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c UserClient) ListByLabCompleteMatchingPredicate(ctx context.Context, id LabId, predicate UserOperationPredicate) (resp ListByLabCompleteResult, err error) {
	items := make([]User, 0)

	page, err := c.ListByLab(ctx, id)
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

	out := ListByLabCompleteResult{
		Items: items,
	}
	return out, nil
}
