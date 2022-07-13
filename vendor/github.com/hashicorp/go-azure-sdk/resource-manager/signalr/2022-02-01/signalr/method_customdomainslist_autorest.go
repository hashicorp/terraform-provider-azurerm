package signalr

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

type CustomDomainsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]CustomDomain

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (CustomDomainsListOperationResponse, error)
}

type CustomDomainsListCompleteResult struct {
	Items []CustomDomain
}

func (r CustomDomainsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r CustomDomainsListOperationResponse) LoadMore(ctx context.Context) (resp CustomDomainsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// CustomDomainsList ...
func (c SignalRClient) CustomDomainsList(ctx context.Context, id SignalRId) (resp CustomDomainsListOperationResponse, err error) {
	req, err := c.preparerForCustomDomainsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomDomainsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomDomainsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForCustomDomainsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomDomainsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// CustomDomainsListComplete retrieves all of the results into a single object
func (c SignalRClient) CustomDomainsListComplete(ctx context.Context, id SignalRId) (CustomDomainsListCompleteResult, error) {
	return c.CustomDomainsListCompleteMatchingPredicate(ctx, id, CustomDomainOperationPredicate{})
}

// CustomDomainsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c SignalRClient) CustomDomainsListCompleteMatchingPredicate(ctx context.Context, id SignalRId, predicate CustomDomainOperationPredicate) (resp CustomDomainsListCompleteResult, err error) {
	items := make([]CustomDomain, 0)

	page, err := c.CustomDomainsList(ctx, id)
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

	out := CustomDomainsListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForCustomDomainsList prepares the CustomDomainsList request.
func (c SignalRClient) preparerForCustomDomainsList(ctx context.Context, id SignalRId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/customDomains", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForCustomDomainsListWithNextLink prepares the CustomDomainsList request with the given nextLink token.
func (c SignalRClient) preparerForCustomDomainsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForCustomDomainsList handles the response to the CustomDomainsList request. The method always
// closes the http.Response Body.
func (c SignalRClient) responderForCustomDomainsList(resp *http.Response) (result CustomDomainsListOperationResponse, err error) {
	type page struct {
		Values   []CustomDomain `json:"value"`
		NextLink *string        `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result CustomDomainsListOperationResponse, err error) {
			req, err := c.preparerForCustomDomainsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomDomainsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomDomainsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForCustomDomainsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "CustomDomainsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
