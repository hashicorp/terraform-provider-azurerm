package signalr

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PrivateLinkResourcesListResponse struct {
	HttpResponse *http.Response
	Model        *[]PrivateLinkResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (PrivateLinkResourcesListResponse, error)
}

type PrivateLinkResourcesListCompleteResult struct {
	Items []PrivateLinkResource
}

func (r PrivateLinkResourcesListResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r PrivateLinkResourcesListResponse) LoadMore(ctx context.Context) (resp PrivateLinkResourcesListResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// PrivateLinkResourcesList ...
func (c SignalRClient) PrivateLinkResourcesList(ctx context.Context, id SignalRId) (resp PrivateLinkResourcesListResponse, err error) {
	req, err := c.preparerForPrivateLinkResourcesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "PrivateLinkResourcesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "PrivateLinkResourcesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForPrivateLinkResourcesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "PrivateLinkResourcesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// PrivateLinkResourcesListComplete retrieves all of the results into a single object
func (c SignalRClient) PrivateLinkResourcesListComplete(ctx context.Context, id SignalRId) (PrivateLinkResourcesListCompleteResult, error) {
	return c.PrivateLinkResourcesListCompleteMatchingPredicate(ctx, id, PrivateLinkResourcePredicate{})
}

// PrivateLinkResourcesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c SignalRClient) PrivateLinkResourcesListCompleteMatchingPredicate(ctx context.Context, id SignalRId, predicate PrivateLinkResourcePredicate) (resp PrivateLinkResourcesListCompleteResult, err error) {
	items := make([]PrivateLinkResource, 0)

	page, err := c.PrivateLinkResourcesList(ctx, id)
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

	out := PrivateLinkResourcesListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForPrivateLinkResourcesList prepares the PrivateLinkResourcesList request.
func (c SignalRClient) preparerForPrivateLinkResourcesList(ctx context.Context, id SignalRId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/privateLinkResources", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForPrivateLinkResourcesListWithNextLink prepares the PrivateLinkResourcesList request with the given nextLink token.
func (c SignalRClient) preparerForPrivateLinkResourcesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForPrivateLinkResourcesList handles the response to the PrivateLinkResourcesList request. The method always
// closes the http.Response Body.
func (c SignalRClient) responderForPrivateLinkResourcesList(resp *http.Response) (result PrivateLinkResourcesListResponse, err error) {
	type page struct {
		Values   []PrivateLinkResource `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result PrivateLinkResourcesListResponse, err error) {
			req, err := c.preparerForPrivateLinkResourcesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "PrivateLinkResourcesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "PrivateLinkResourcesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForPrivateLinkResourcesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "signalr.SignalRClient", "PrivateLinkResourcesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
