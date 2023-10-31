package tableservice

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

type TableListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Table

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (TableListOperationResponse, error)
}

type TableListCompleteResult struct {
	Items []Table
}

func (r TableListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r TableListOperationResponse) LoadMore(ctx context.Context) (resp TableListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// TableList ...
func (c TableServiceClient) TableList(ctx context.Context, id commonids.StorageAccountId) (resp TableListOperationResponse, err error) {
	req, err := c.preparerForTableList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForTableList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForTableList prepares the TableList request.
func (c TableServiceClient) preparerForTableList(ctx context.Context, id commonids.StorageAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/tableServices/default/tables", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForTableListWithNextLink prepares the TableList request with the given nextLink token.
func (c TableServiceClient) preparerForTableListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForTableList handles the response to the TableList request. The method always
// closes the http.Response Body.
func (c TableServiceClient) responderForTableList(resp *http.Response) (result TableListOperationResponse, err error) {
	type page struct {
		Values   []Table `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result TableListOperationResponse, err error) {
			req, err := c.preparerForTableListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForTableList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// TableListComplete retrieves all of the results into a single object
func (c TableServiceClient) TableListComplete(ctx context.Context, id commonids.StorageAccountId) (TableListCompleteResult, error) {
	return c.TableListCompleteMatchingPredicate(ctx, id, TableOperationPredicate{})
}

// TableListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c TableServiceClient) TableListCompleteMatchingPredicate(ctx context.Context, id commonids.StorageAccountId, predicate TableOperationPredicate) (resp TableListCompleteResult, err error) {
	items := make([]Table, 0)

	page, err := c.TableList(ctx, id)
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

	out := TableListCompleteResult{
		Items: items,
	}
	return out, nil
}
