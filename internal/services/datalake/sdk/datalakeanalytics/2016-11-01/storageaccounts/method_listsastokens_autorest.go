package storageaccounts

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListSasTokensResponse struct {
	HttpResponse *http.Response
	Model        *[]SasTokenInformation

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListSasTokensResponse, error)
}

type ListSasTokensCompleteResult struct {
	Items []SasTokenInformation
}

func (r ListSasTokensResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListSasTokensResponse) LoadMore(ctx context.Context) (resp ListSasTokensResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListSasTokens ...
func (c StorageAccountsClient) ListSasTokens(ctx context.Context, id ContainerId) (resp ListSasTokensResponse, err error) {
	req, err := c.preparerForListSasTokens(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListSasTokens", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListSasTokens", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListSasTokens(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListSasTokens", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListSasTokensComplete retrieves all of the results into a single object
func (c StorageAccountsClient) ListSasTokensComplete(ctx context.Context, id ContainerId) (ListSasTokensCompleteResult, error) {
	return c.ListSasTokensCompleteMatchingPredicate(ctx, id, SasTokenInformationPredicate{})
}

// ListSasTokensCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c StorageAccountsClient) ListSasTokensCompleteMatchingPredicate(ctx context.Context, id ContainerId, predicate SasTokenInformationPredicate) (resp ListSasTokensCompleteResult, err error) {
	items := make([]SasTokenInformation, 0)

	page, err := c.ListSasTokens(ctx, id)
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

	out := ListSasTokensCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListSasTokens prepares the ListSasTokens request.
func (c StorageAccountsClient) preparerForListSasTokens(ctx context.Context, id ContainerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listSasTokens", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListSasTokensWithNextLink prepares the ListSasTokens request with the given nextLink token.
func (c StorageAccountsClient) preparerForListSasTokensWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListSasTokens handles the response to the ListSasTokens request. The method always
// closes the http.Response Body.
func (c StorageAccountsClient) responderForListSasTokens(resp *http.Response) (result ListSasTokensResponse, err error) {
	type page struct {
		Values   []SasTokenInformation `json:"value"`
		NextLink *string               `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListSasTokensResponse, err error) {
			req, err := c.preparerForListSasTokensWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListSasTokens", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListSasTokens", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListSasTokens(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListSasTokens", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
