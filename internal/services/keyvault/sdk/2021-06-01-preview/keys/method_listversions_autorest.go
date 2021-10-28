package keys

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListVersionsResponse struct {
	HttpResponse *http.Response
	Model        *[]Key

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListVersionsResponse, error)
}

type ListVersionsCompleteResult struct {
	Items []Key
}

func (r ListVersionsResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListVersionsResponse) LoadMore(ctx context.Context) (resp ListVersionsResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListVersions ...
func (c KeysClient) ListVersions(ctx context.Context, id KeyId) (resp ListVersionsResponse, err error) {
	req, err := c.preparerForListVersions(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "keys.KeysClient", "ListVersions", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "keys.KeysClient", "ListVersions", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListVersions(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "keys.KeysClient", "ListVersions", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListVersionsComplete retrieves all of the results into a single object
func (c KeysClient) ListVersionsComplete(ctx context.Context, id KeyId) (ListVersionsCompleteResult, error) {
	return c.ListVersionsCompleteMatchingPredicate(ctx, id, KeyPredicate{})
}

// ListVersionsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c KeysClient) ListVersionsCompleteMatchingPredicate(ctx context.Context, id KeyId, predicate KeyPredicate) (resp ListVersionsCompleteResult, err error) {
	items := make([]Key, 0)

	page, err := c.ListVersions(ctx, id)
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

	out := ListVersionsCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListVersions prepares the ListVersions request.
func (c KeysClient) preparerForListVersions(ctx context.Context, id KeyId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/versions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListVersionsWithNextLink prepares the ListVersions request with the given nextLink token.
func (c KeysClient) preparerForListVersionsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListVersions handles the response to the ListVersions request. The method always
// closes the http.Response Body.
func (c KeysClient) responderForListVersions(resp *http.Response) (result ListVersionsResponse, err error) {
	type page struct {
		Values   []Key   `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListVersionsResponse, err error) {
			req, err := c.preparerForListVersionsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "keys.KeysClient", "ListVersions", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "keys.KeysClient", "ListVersions", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListVersions(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "keys.KeysClient", "ListVersions", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
