package wcfrelays

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByNamespaceResponse struct {
	HttpResponse *http.Response
	Model        *[]WcfRelay

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByNamespaceResponse, error)
}

type ListByNamespaceCompleteResult struct {
	Items []WcfRelay
}

func (r ListByNamespaceResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByNamespaceResponse) LoadMore(ctx context.Context) (resp ListByNamespaceResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByNamespace ...
func (c WCFRelaysClient) ListByNamespace(ctx context.Context, id NamespaceId) (resp ListByNamespaceResponse, err error) {
	req, err := c.preparerForListByNamespace(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "wcfrelays.WCFRelaysClient", "ListByNamespace", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "wcfrelays.WCFRelaysClient", "ListByNamespace", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByNamespace(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "wcfrelays.WCFRelaysClient", "ListByNamespace", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByNamespaceComplete retrieves all of the results into a single object
func (c WCFRelaysClient) ListByNamespaceComplete(ctx context.Context, id NamespaceId) (ListByNamespaceCompleteResult, error) {
	return c.ListByNamespaceCompleteMatchingPredicate(ctx, id, WcfRelayPredicate{})
}

// ListByNamespaceCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c WCFRelaysClient) ListByNamespaceCompleteMatchingPredicate(ctx context.Context, id NamespaceId, predicate WcfRelayPredicate) (resp ListByNamespaceCompleteResult, err error) {
	items := make([]WcfRelay, 0)

	page, err := c.ListByNamespace(ctx, id)
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

	out := ListByNamespaceCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByNamespace prepares the ListByNamespace request.
func (c WCFRelaysClient) preparerForListByNamespace(ctx context.Context, id NamespaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/wcfRelays", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByNamespaceWithNextLink prepares the ListByNamespace request with the given nextLink token.
func (c WCFRelaysClient) preparerForListByNamespaceWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByNamespace handles the response to the ListByNamespace request. The method always
// closes the http.Response Body.
func (c WCFRelaysClient) responderForListByNamespace(resp *http.Response) (result ListByNamespaceResponse, err error) {
	type page struct {
		Values   []WcfRelay `json:"value"`
		NextLink *string    `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByNamespaceResponse, err error) {
			req, err := c.preparerForListByNamespaceWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "wcfrelays.WCFRelaysClient", "ListByNamespace", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "wcfrelays.WCFRelaysClient", "ListByNamespace", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByNamespace(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "wcfrelays.WCFRelaysClient", "ListByNamespace", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
