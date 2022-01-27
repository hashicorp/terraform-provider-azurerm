package afdorigins

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByOriginGroupResponse struct {
	HttpResponse *http.Response
	Model        *[]AFDOrigin

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByOriginGroupResponse, error)
}

type ListByOriginGroupCompleteResult struct {
	Items []AFDOrigin
}

func (r ListByOriginGroupResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByOriginGroupResponse) LoadMore(ctx context.Context) (resp ListByOriginGroupResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByOriginGroup ...
func (c AFDOriginsClient) ListByOriginGroup(ctx context.Context, id OriginGroupId) (resp ListByOriginGroupResponse, err error) {
	req, err := c.preparerForListByOriginGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "afdorigins.AFDOriginsClient", "ListByOriginGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "afdorigins.AFDOriginsClient", "ListByOriginGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByOriginGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "afdorigins.AFDOriginsClient", "ListByOriginGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByOriginGroupComplete retrieves all of the results into a single object
func (c AFDOriginsClient) ListByOriginGroupComplete(ctx context.Context, id OriginGroupId) (ListByOriginGroupCompleteResult, error) {
	return c.ListByOriginGroupCompleteMatchingPredicate(ctx, id, AFDOriginPredicate{})
}

// ListByOriginGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AFDOriginsClient) ListByOriginGroupCompleteMatchingPredicate(ctx context.Context, id OriginGroupId, predicate AFDOriginPredicate) (resp ListByOriginGroupCompleteResult, err error) {
	items := make([]AFDOrigin, 0)

	page, err := c.ListByOriginGroup(ctx, id)
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

	out := ListByOriginGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByOriginGroup prepares the ListByOriginGroup request.
func (c AFDOriginsClient) preparerForListByOriginGroup(ctx context.Context, id OriginGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/origins", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByOriginGroupWithNextLink prepares the ListByOriginGroup request with the given nextLink token.
func (c AFDOriginsClient) preparerForListByOriginGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByOriginGroup handles the response to the ListByOriginGroup request. The method always
// closes the http.Response Body.
func (c AFDOriginsClient) responderForListByOriginGroup(resp *http.Response) (result ListByOriginGroupResponse, err error) {
	type page struct {
		Values   []AFDOrigin `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByOriginGroupResponse, err error) {
			req, err := c.preparerForListByOriginGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "afdorigins.AFDOriginsClient", "ListByOriginGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "afdorigins.AFDOriginsClient", "ListByOriginGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByOriginGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "afdorigins.AFDOriginsClient", "ListByOriginGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
