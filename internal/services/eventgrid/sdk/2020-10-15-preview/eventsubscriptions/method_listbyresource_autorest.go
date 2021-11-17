package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByResourceResponse struct {
	HttpResponse *http.Response
	Model        *[]EventSubscription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByResourceResponse, error)
}

type ListByResourceCompleteResult struct {
	Items []EventSubscription
}

func (r ListByResourceResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByResourceResponse) LoadMore(ctx context.Context) (resp ListByResourceResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByResourceOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListByResourceOptions() ListByResourceOptions {
	return ListByResourceOptions{}
}

func (o ListByResourceOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListByResource ...
func (c EventSubscriptionsClient) ListByResource(ctx context.Context, id ProviderId, options ListByResourceOptions) (resp ListByResourceResponse, err error) {
	req, err := c.preparerForListByResource(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListByResource", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListByResource", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByResource(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListByResource", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByResourceComplete retrieves all of the results into a single object
func (c EventSubscriptionsClient) ListByResourceComplete(ctx context.Context, id ProviderId, options ListByResourceOptions) (ListByResourceCompleteResult, error) {
	return c.ListByResourceCompleteMatchingPredicate(ctx, id, options, EventSubscriptionPredicate{})
}

// ListByResourceCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EventSubscriptionsClient) ListByResourceCompleteMatchingPredicate(ctx context.Context, id ProviderId, options ListByResourceOptions, predicate EventSubscriptionPredicate) (resp ListByResourceCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	page, err := c.ListByResource(ctx, id, options)
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

	out := ListByResourceCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByResource prepares the ListByResource request.
func (c EventSubscriptionsClient) preparerForListByResource(ctx context.Context, id ProviderId, options ListByResourceOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.EventGrid/eventSubscriptions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByResourceWithNextLink prepares the ListByResource request with the given nextLink token.
func (c EventSubscriptionsClient) preparerForListByResourceWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByResource handles the response to the ListByResource request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForListByResource(resp *http.Response) (result ListByResourceResponse, err error) {
	type page struct {
		Values   []EventSubscription `json:"value"`
		NextLink *string             `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByResourceResponse, err error) {
			req, err := c.preparerForListByResourceWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListByResource", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListByResource", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByResource(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListByResource", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
