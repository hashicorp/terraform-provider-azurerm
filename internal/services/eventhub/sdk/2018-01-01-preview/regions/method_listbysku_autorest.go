package regions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListBySkuResponse struct {
	HttpResponse *http.Response
	Model        *[]MessagingRegions

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListBySkuResponse, error)
}

type ListBySkuCompleteResult struct {
	Items []MessagingRegions
}

func (r ListBySkuResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListBySkuResponse) LoadMore(ctx context.Context) (resp ListBySkuResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListBySku ...
func (c RegionsClient) ListBySku(ctx context.Context, id SkuId) (resp ListBySkuResponse, err error) {
	req, err := c.preparerForListBySku(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "regions.RegionsClient", "ListBySku", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "regions.RegionsClient", "ListBySku", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListBySku(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "regions.RegionsClient", "ListBySku", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListBySkuComplete retrieves all of the results into a single object
func (c RegionsClient) ListBySkuComplete(ctx context.Context, id SkuId) (ListBySkuCompleteResult, error) {
	return c.ListBySkuCompleteMatchingPredicate(ctx, id, MessagingRegionsPredicate{})
}

// ListBySkuCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c RegionsClient) ListBySkuCompleteMatchingPredicate(ctx context.Context, id SkuId, predicate MessagingRegionsPredicate) (resp ListBySkuCompleteResult, err error) {
	items := make([]MessagingRegions, 0)

	page, err := c.ListBySku(ctx, id)
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

	out := ListBySkuCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListBySku prepares the ListBySku request.
func (c RegionsClient) preparerForListBySku(ctx context.Context, id SkuId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/regions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListBySkuWithNextLink prepares the ListBySku request with the given nextLink token.
func (c RegionsClient) preparerForListBySkuWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListBySku handles the response to the ListBySku request. The method always
// closes the http.Response Body.
func (c RegionsClient) responderForListBySku(resp *http.Response) (result ListBySkuResponse, err error) {
	type page struct {
		Values   []MessagingRegions `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListBySkuResponse, err error) {
			req, err := c.preparerForListBySkuWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "regions.RegionsClient", "ListBySku", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "regions.RegionsClient", "ListBySku", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListBySku(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "regions.RegionsClient", "ListBySku", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
