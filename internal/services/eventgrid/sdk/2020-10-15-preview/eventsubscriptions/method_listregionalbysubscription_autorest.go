package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListRegionalBySubscriptionResponse struct {
	HttpResponse *http.Response
	Model        *[]EventSubscription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListRegionalBySubscriptionResponse, error)
}

type ListRegionalBySubscriptionCompleteResult struct {
	Items []EventSubscription
}

func (r ListRegionalBySubscriptionResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListRegionalBySubscriptionResponse) LoadMore(ctx context.Context) (resp ListRegionalBySubscriptionResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListRegionalBySubscriptionOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListRegionalBySubscriptionOptions() ListRegionalBySubscriptionOptions {
	return ListRegionalBySubscriptionOptions{}
}

func (o ListRegionalBySubscriptionOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListRegionalBySubscription ...
func (c EventSubscriptionsClient) ListRegionalBySubscription(ctx context.Context, id LocationId, options ListRegionalBySubscriptionOptions) (resp ListRegionalBySubscriptionResponse, err error) {
	req, err := c.preparerForListRegionalBySubscription(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListRegionalBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListRegionalBySubscriptionComplete retrieves all of the results into a single object
func (c EventSubscriptionsClient) ListRegionalBySubscriptionComplete(ctx context.Context, id LocationId, options ListRegionalBySubscriptionOptions) (ListRegionalBySubscriptionCompleteResult, error) {
	return c.ListRegionalBySubscriptionCompleteMatchingPredicate(ctx, id, options, EventSubscriptionPredicate{})
}

// ListRegionalBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EventSubscriptionsClient) ListRegionalBySubscriptionCompleteMatchingPredicate(ctx context.Context, id LocationId, options ListRegionalBySubscriptionOptions, predicate EventSubscriptionPredicate) (resp ListRegionalBySubscriptionCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	page, err := c.ListRegionalBySubscription(ctx, id, options)
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

	out := ListRegionalBySubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListRegionalBySubscription prepares the ListRegionalBySubscription request.
func (c EventSubscriptionsClient) preparerForListRegionalBySubscription(ctx context.Context, id LocationId, options ListRegionalBySubscriptionOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/eventSubscriptions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListRegionalBySubscriptionWithNextLink prepares the ListRegionalBySubscription request with the given nextLink token.
func (c EventSubscriptionsClient) preparerForListRegionalBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListRegionalBySubscription handles the response to the ListRegionalBySubscription request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForListRegionalBySubscription(resp *http.Response) (result ListRegionalBySubscriptionResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListRegionalBySubscriptionResponse, err error) {
			req, err := c.preparerForListRegionalBySubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListRegionalBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
