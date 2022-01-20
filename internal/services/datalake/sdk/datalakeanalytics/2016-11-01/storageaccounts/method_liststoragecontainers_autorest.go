package storageaccounts

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListStorageContainersResponse struct {
	HttpResponse *http.Response
	Model        *[]StorageContainer

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListStorageContainersResponse, error)
}

type ListStorageContainersCompleteResult struct {
	Items []StorageContainer
}

func (r ListStorageContainersResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListStorageContainersResponse) LoadMore(ctx context.Context) (resp ListStorageContainersResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListStorageContainers ...
func (c StorageAccountsClient) ListStorageContainers(ctx context.Context, id StorageAccountId) (resp ListStorageContainersResponse, err error) {
	req, err := c.preparerForListStorageContainers(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListStorageContainers", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListStorageContainers", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListStorageContainers(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListStorageContainers", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListStorageContainersComplete retrieves all of the results into a single object
func (c StorageAccountsClient) ListStorageContainersComplete(ctx context.Context, id StorageAccountId) (ListStorageContainersCompleteResult, error) {
	return c.ListStorageContainersCompleteMatchingPredicate(ctx, id, StorageContainerPredicate{})
}

// ListStorageContainersCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c StorageAccountsClient) ListStorageContainersCompleteMatchingPredicate(ctx context.Context, id StorageAccountId, predicate StorageContainerPredicate) (resp ListStorageContainersCompleteResult, err error) {
	items := make([]StorageContainer, 0)

	page, err := c.ListStorageContainers(ctx, id)
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

	out := ListStorageContainersCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListStorageContainers prepares the ListStorageContainers request.
func (c StorageAccountsClient) preparerForListStorageContainers(ctx context.Context, id StorageAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/containers", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListStorageContainersWithNextLink prepares the ListStorageContainers request with the given nextLink token.
func (c StorageAccountsClient) preparerForListStorageContainersWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListStorageContainers handles the response to the ListStorageContainers request. The method always
// closes the http.Response Body.
func (c StorageAccountsClient) responderForListStorageContainers(resp *http.Response) (result ListStorageContainersResponse, err error) {
	type page struct {
		Values   []StorageContainer `json:"value"`
		NextLink *string            `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListStorageContainersResponse, err error) {
			req, err := c.preparerForListStorageContainersWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListStorageContainers", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListStorageContainers", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListStorageContainers(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "storageaccounts.StorageAccountsClient", "ListStorageContainers", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
