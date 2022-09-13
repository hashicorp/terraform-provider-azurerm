package monitors

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

type ListLinkableEnvironmentsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]LinkableEnvironmentResponse

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListLinkableEnvironmentsOperationResponse, error)
}

type ListLinkableEnvironmentsCompleteResult struct {
	Items []LinkableEnvironmentResponse
}

func (r ListLinkableEnvironmentsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListLinkableEnvironmentsOperationResponse) LoadMore(ctx context.Context) (resp ListLinkableEnvironmentsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListLinkableEnvironments ...
func (c MonitorsClient) ListLinkableEnvironments(ctx context.Context, id MonitorId, input LinkableEnvironmentRequest) (resp ListLinkableEnvironmentsOperationResponse, err error) {
	req, err := c.preparerForListLinkableEnvironments(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListLinkableEnvironments", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListLinkableEnvironments", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListLinkableEnvironments(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListLinkableEnvironments", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListLinkableEnvironments prepares the ListLinkableEnvironments request.
func (c MonitorsClient) preparerForListLinkableEnvironments(ctx context.Context, id MonitorId, input LinkableEnvironmentRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listLinkableEnvironments", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListLinkableEnvironmentsWithNextLink prepares the ListLinkableEnvironments request with the given nextLink token.
func (c MonitorsClient) preparerForListLinkableEnvironmentsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListLinkableEnvironments handles the response to the ListLinkableEnvironments request. The method always
// closes the http.Response Body.
func (c MonitorsClient) responderForListLinkableEnvironments(resp *http.Response) (result ListLinkableEnvironmentsOperationResponse, err error) {
	type page struct {
		Values   []LinkableEnvironmentResponse `json:"value"`
		NextLink *string                       `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListLinkableEnvironmentsOperationResponse, err error) {
			req, err := c.preparerForListLinkableEnvironmentsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListLinkableEnvironments", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListLinkableEnvironments", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListLinkableEnvironments(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListLinkableEnvironments", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListLinkableEnvironmentsComplete retrieves all of the results into a single object
func (c MonitorsClient) ListLinkableEnvironmentsComplete(ctx context.Context, id MonitorId, input LinkableEnvironmentRequest) (ListLinkableEnvironmentsCompleteResult, error) {
	return c.ListLinkableEnvironmentsCompleteMatchingPredicate(ctx, id, input, LinkableEnvironmentResponseOperationPredicate{})
}

// ListLinkableEnvironmentsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c MonitorsClient) ListLinkableEnvironmentsCompleteMatchingPredicate(ctx context.Context, id MonitorId, input LinkableEnvironmentRequest, predicate LinkableEnvironmentResponseOperationPredicate) (resp ListLinkableEnvironmentsCompleteResult, err error) {
	items := make([]LinkableEnvironmentResponse, 0)

	page, err := c.ListLinkableEnvironments(ctx, id, input)
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

	out := ListLinkableEnvironmentsCompleteResult{
		Items: items,
	}
	return out, nil
}
