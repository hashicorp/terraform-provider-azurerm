package servicelinker

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkerListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]LinkerResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (LinkerListOperationResponse, error)
}

type LinkerListCompleteResult struct {
	Items []LinkerResource
}

func (r LinkerListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r LinkerListOperationResponse) LoadMore(ctx context.Context) (resp LinkerListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// LinkerList ...
func (c ServiceLinkerClient) LinkerList(ctx context.Context, id commonids.ScopeId) (resp LinkerListOperationResponse, err error) {
	req, err := c.preparerForLinkerList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicelinker.ServiceLinkerClient", "LinkerList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicelinker.ServiceLinkerClient", "LinkerList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForLinkerList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicelinker.ServiceLinkerClient", "LinkerList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForLinkerList prepares the LinkerList request.
func (c ServiceLinkerClient) preparerForLinkerList(ctx context.Context, id commonids.ScopeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.ServiceLinker/linkers", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForLinkerListWithNextLink prepares the LinkerList request with the given nextLink token.
func (c ServiceLinkerClient) preparerForLinkerListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForLinkerList handles the response to the LinkerList request. The method always
// closes the http.Response Body.
func (c ServiceLinkerClient) responderForLinkerList(resp *http.Response) (result LinkerListOperationResponse, err error) {
	type page struct {
		Values   []LinkerResource `json:"value"`
		NextLink *string          `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result LinkerListOperationResponse, err error) {
			req, err := c.preparerForLinkerListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "servicelinker.ServiceLinkerClient", "LinkerList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "servicelinker.ServiceLinkerClient", "LinkerList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForLinkerList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "servicelinker.ServiceLinkerClient", "LinkerList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// LinkerListComplete retrieves all of the results into a single object
func (c ServiceLinkerClient) LinkerListComplete(ctx context.Context, id commonids.ScopeId) (LinkerListCompleteResult, error) {
	return c.LinkerListCompleteMatchingPredicate(ctx, id, LinkerResourceOperationPredicate{})
}

// LinkerListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ServiceLinkerClient) LinkerListCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, predicate LinkerResourceOperationPredicate) (resp LinkerListCompleteResult, err error) {
	items := make([]LinkerResource, 0)

	page, err := c.LinkerList(ctx, id)
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

	out := LinkerListCompleteResult{
		Items: items,
	}
	return out, nil
}
