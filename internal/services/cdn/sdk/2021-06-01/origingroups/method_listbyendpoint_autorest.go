package origingroups

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByEndpointResponse struct {
	HttpResponse *http.Response
	Model        *[]OriginGroup

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByEndpointResponse, error)
}

type ListByEndpointCompleteResult struct {
	Items []OriginGroup
}

func (r ListByEndpointResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByEndpointResponse) LoadMore(ctx context.Context) (resp ListByEndpointResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByEndpoint ...
func (c OriginGroupsClient) ListByEndpoint(ctx context.Context, id EndpointId) (resp ListByEndpointResponse, err error) {
	req, err := c.preparerForListByEndpoint(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "origingroups.OriginGroupsClient", "ListByEndpoint", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "origingroups.OriginGroupsClient", "ListByEndpoint", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByEndpoint(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "origingroups.OriginGroupsClient", "ListByEndpoint", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByEndpointComplete retrieves all of the results into a single object
func (c OriginGroupsClient) ListByEndpointComplete(ctx context.Context, id EndpointId) (ListByEndpointCompleteResult, error) {
	return c.ListByEndpointCompleteMatchingPredicate(ctx, id, OriginGroupPredicate{})
}

// ListByEndpointCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c OriginGroupsClient) ListByEndpointCompleteMatchingPredicate(ctx context.Context, id EndpointId, predicate OriginGroupPredicate) (resp ListByEndpointCompleteResult, err error) {
	items := make([]OriginGroup, 0)

	page, err := c.ListByEndpoint(ctx, id)
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

	out := ListByEndpointCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByEndpoint prepares the ListByEndpoint request.
func (c OriginGroupsClient) preparerForListByEndpoint(ctx context.Context, id EndpointId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/originGroups", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByEndpointWithNextLink prepares the ListByEndpoint request with the given nextLink token.
func (c OriginGroupsClient) preparerForListByEndpointWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByEndpoint handles the response to the ListByEndpoint request. The method always
// closes the http.Response Body.
func (c OriginGroupsClient) responderForListByEndpoint(resp *http.Response) (result ListByEndpointResponse, err error) {
	type page struct {
		Values   []OriginGroup `json:"value"`
		NextLink *string       `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByEndpointResponse, err error) {
			req, err := c.preparerForListByEndpointWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "origingroups.OriginGroupsClient", "ListByEndpoint", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "origingroups.OriginGroupsClient", "ListByEndpoint", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByEndpoint(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "origingroups.OriginGroupsClient", "ListByEndpoint", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
