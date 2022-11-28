package accounts

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

type MediaservicesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]MediaService

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (MediaservicesListOperationResponse, error)
}

type MediaservicesListCompleteResult struct {
	Items []MediaService
}

func (r MediaservicesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r MediaservicesListOperationResponse) LoadMore(ctx context.Context) (resp MediaservicesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// MediaservicesList ...
func (c AccountsClient) MediaservicesList(ctx context.Context, id commonids.ResourceGroupId) (resp MediaservicesListOperationResponse, err error) {
	req, err := c.preparerForMediaservicesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForMediaservicesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForMediaservicesList prepares the MediaservicesList request.
func (c AccountsClient) preparerForMediaservicesList(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Media/mediaServices", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForMediaservicesListWithNextLink prepares the MediaservicesList request with the given nextLink token.
func (c AccountsClient) preparerForMediaservicesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForMediaservicesList handles the response to the MediaservicesList request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForMediaservicesList(resp *http.Response) (result MediaservicesListOperationResponse, err error) {
	type page struct {
		Values   []MediaService `json:"value"`
		NextLink *string        `json:"@odata.nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result MediaservicesListOperationResponse, err error) {
			req, err := c.preparerForMediaservicesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForMediaservicesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// MediaservicesListComplete retrieves all of the results into a single object
func (c AccountsClient) MediaservicesListComplete(ctx context.Context, id commonids.ResourceGroupId) (MediaservicesListCompleteResult, error) {
	return c.MediaservicesListCompleteMatchingPredicate(ctx, id, MediaServiceOperationPredicate{})
}

// MediaservicesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AccountsClient) MediaservicesListCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate MediaServiceOperationPredicate) (resp MediaservicesListCompleteResult, err error) {
	items := make([]MediaService, 0)

	page, err := c.MediaservicesList(ctx, id)
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

	out := MediaservicesListCompleteResult{
		Items: items,
	}
	return out, nil
}
