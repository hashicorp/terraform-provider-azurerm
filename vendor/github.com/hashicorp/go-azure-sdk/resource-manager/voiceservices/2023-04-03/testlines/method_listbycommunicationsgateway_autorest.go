package testlines

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

type ListByCommunicationsGatewayOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]TestLine

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByCommunicationsGatewayOperationResponse, error)
}

type ListByCommunicationsGatewayCompleteResult struct {
	Items []TestLine
}

func (r ListByCommunicationsGatewayOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByCommunicationsGatewayOperationResponse) LoadMore(ctx context.Context) (resp ListByCommunicationsGatewayOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByCommunicationsGateway ...
func (c TestLinesClient) ListByCommunicationsGateway(ctx context.Context, id CommunicationsGatewayId) (resp ListByCommunicationsGatewayOperationResponse, err error) {
	req, err := c.preparerForListByCommunicationsGateway(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "testlines.TestLinesClient", "ListByCommunicationsGateway", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "testlines.TestLinesClient", "ListByCommunicationsGateway", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByCommunicationsGateway(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "testlines.TestLinesClient", "ListByCommunicationsGateway", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByCommunicationsGateway prepares the ListByCommunicationsGateway request.
func (c TestLinesClient) preparerForListByCommunicationsGateway(ctx context.Context, id CommunicationsGatewayId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/testLines", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByCommunicationsGatewayWithNextLink prepares the ListByCommunicationsGateway request with the given nextLink token.
func (c TestLinesClient) preparerForListByCommunicationsGatewayWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByCommunicationsGateway handles the response to the ListByCommunicationsGateway request. The method always
// closes the http.Response Body.
func (c TestLinesClient) responderForListByCommunicationsGateway(resp *http.Response) (result ListByCommunicationsGatewayOperationResponse, err error) {
	type page struct {
		Values   []TestLine `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByCommunicationsGatewayOperationResponse, err error) {
			req, err := c.preparerForListByCommunicationsGatewayWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "testlines.TestLinesClient", "ListByCommunicationsGateway", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "testlines.TestLinesClient", "ListByCommunicationsGateway", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByCommunicationsGateway(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "testlines.TestLinesClient", "ListByCommunicationsGateway", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByCommunicationsGatewayComplete retrieves all of the results into a single object
func (c TestLinesClient) ListByCommunicationsGatewayComplete(ctx context.Context, id CommunicationsGatewayId) (ListByCommunicationsGatewayCompleteResult, error) {
	return c.ListByCommunicationsGatewayCompleteMatchingPredicate(ctx, id, TestLineOperationPredicate{})
}

// ListByCommunicationsGatewayCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c TestLinesClient) ListByCommunicationsGatewayCompleteMatchingPredicate(ctx context.Context, id CommunicationsGatewayId, predicate TestLineOperationPredicate) (resp ListByCommunicationsGatewayCompleteResult, err error) {
	items := make([]TestLine, 0)

	page, err := c.ListByCommunicationsGateway(ctx, id)
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

	out := ListByCommunicationsGatewayCompleteResult{
		Items: items,
	}
	return out, nil
}
