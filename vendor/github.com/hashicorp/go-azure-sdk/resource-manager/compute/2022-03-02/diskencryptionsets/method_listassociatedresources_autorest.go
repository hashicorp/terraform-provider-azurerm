package diskencryptionsets

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

type ListAssociatedResourcesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]string

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListAssociatedResourcesOperationResponse, error)
}

type ListAssociatedResourcesCompleteResult struct {
	Items []string
}

func (r ListAssociatedResourcesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListAssociatedResourcesOperationResponse) LoadMore(ctx context.Context) (resp ListAssociatedResourcesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListAssociatedResources ...
func (c DiskEncryptionSetsClient) ListAssociatedResources(ctx context.Context, id DiskEncryptionSetId) (resp ListAssociatedResourcesOperationResponse, err error) {
	req, err := c.preparerForListAssociatedResources(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskencryptionsets.DiskEncryptionSetsClient", "ListAssociatedResources", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskencryptionsets.DiskEncryptionSetsClient", "ListAssociatedResources", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListAssociatedResources(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "diskencryptionsets.DiskEncryptionSetsClient", "ListAssociatedResources", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListAssociatedResources prepares the ListAssociatedResources request.
func (c DiskEncryptionSetsClient) preparerForListAssociatedResources(ctx context.Context, id DiskEncryptionSetId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/associatedResources", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListAssociatedResourcesWithNextLink prepares the ListAssociatedResources request with the given nextLink token.
func (c DiskEncryptionSetsClient) preparerForListAssociatedResourcesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListAssociatedResources handles the response to the ListAssociatedResources request. The method always
// closes the http.Response Body.
func (c DiskEncryptionSetsClient) responderForListAssociatedResources(resp *http.Response) (result ListAssociatedResourcesOperationResponse, err error) {
	type page struct {
		Values   []string `json:"value"`
		NextLink *string  `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListAssociatedResourcesOperationResponse, err error) {
			req, err := c.preparerForListAssociatedResourcesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "diskencryptionsets.DiskEncryptionSetsClient", "ListAssociatedResources", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "diskencryptionsets.DiskEncryptionSetsClient", "ListAssociatedResources", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListAssociatedResources(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "diskencryptionsets.DiskEncryptionSetsClient", "ListAssociatedResources", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListAssociatedResourcesComplete retrieves all of the results into a single object
func (c DiskEncryptionSetsClient) ListAssociatedResourcesComplete(ctx context.Context, id DiskEncryptionSetId) (result ListAssociatedResourcesCompleteResult, err error) {
	items := make([]string, 0)

	page, err := c.ListAssociatedResources(ctx, id)
	if err != nil {
		err = fmt.Errorf("loading the initial page: %+v", err)
		return
	}
	if page.Model != nil {
		for _, v := range *page.Model {
			items = append(items, v)
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
				items = append(items, v)
			}
		}
	}

	out := ListAssociatedResourcesCompleteResult{
		Items: items,
	}
	return out, nil
}
