package configurationstores

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListKeysResponse struct {
	HttpResponse *http.Response
	Model        *[]AccessKey

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListKeysResponse, error)
}

type ListKeysCompleteResult struct {
	Items []AccessKey
}

func (r ListKeysResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListKeysResponse) LoadMore(ctx context.Context) (resp ListKeysResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type AccessKeyPredicate struct {
	// TODO: implement me
}

func (p AccessKeyPredicate) Matches(input AccessKey) bool {
	// TODO: implement me
	// if p.Name != nil && input.Name != *p.Name {
	// 	return false
	// }

	return true
}

// ListKeys ...
func (c ConfigurationStoresClient) ListKeys(ctx context.Context, id ConfigurationStoreId) (resp ListKeysResponse, err error) {
	req, err := c.preparerForListKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationstores.ConfigurationStoresClient", "ListKeys", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationstores.ConfigurationStoresClient", "ListKeys", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListKeys(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurationstores.ConfigurationStoresClient", "ListKeys", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListKeysCompleteMatchingPredicate retrieves all of the results into a single object
func (c ConfigurationStoresClient) ListKeysComplete(ctx context.Context, id ConfigurationStoreId) (ListKeysCompleteResult, error) {
	return c.ListKeysCompleteMatchingPredicate(ctx, id, AccessKeyPredicate{})
}

// ListKeysCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ConfigurationStoresClient) ListKeysCompleteMatchingPredicate(ctx context.Context, id ConfigurationStoreId, predicate AccessKeyPredicate) (resp ListKeysCompleteResult, err error) {
	items := make([]AccessKey, 0)

	page, err := c.ListKeys(ctx, id)
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

	out := ListKeysCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListKeys prepares the ListKeys request.
func (c ConfigurationStoresClient) preparerForListKeys(ctx context.Context, id ConfigurationStoreId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListKeysWithNextLink prepares the ListKeys request with the given nextLink token.
func (c ConfigurationStoresClient) preparerForListKeysWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListKeys handles the response to the ListKeys request. The method always
// closes the http.Response Body.
func (c ConfigurationStoresClient) responderForListKeys(resp *http.Response) (result ListKeysResponse, err error) {
	type page struct {
		Values   []AccessKey `json:"value"`
		NextLink *string     `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListKeysResponse, err error) {
			req, err := c.preparerForListKeysWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "configurationstores.ConfigurationStoresClient", "ListKeys", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "configurationstores.ConfigurationStoresClient", "ListKeys", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListKeys(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "configurationstores.ConfigurationStoresClient", "ListKeys", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
