package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListRegionalBySubscriptionForTopicTypeResponse struct {
	HttpResponse *http.Response
	Model        *[]EventSubscription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListRegionalBySubscriptionForTopicTypeResponse, error)
}

type ListRegionalBySubscriptionForTopicTypeCompleteResult struct {
	Items []EventSubscription
}

func (r ListRegionalBySubscriptionForTopicTypeResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListRegionalBySubscriptionForTopicTypeResponse) LoadMore(ctx context.Context) (resp ListRegionalBySubscriptionForTopicTypeResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListRegionalBySubscriptionForTopicTypeOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListRegionalBySubscriptionForTopicTypeOptions() ListRegionalBySubscriptionForTopicTypeOptions {
	return ListRegionalBySubscriptionForTopicTypeOptions{}
}

func (o ListRegionalBySubscriptionForTopicTypeOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListRegionalBySubscriptionForTopicType ...
func (c EventSubscriptionsClient) ListRegionalBySubscriptionForTopicType(ctx context.Context, id LocationTopicTypeId, options ListRegionalBySubscriptionForTopicTypeOptions) (resp ListRegionalBySubscriptionForTopicTypeResponse, err error) {
	req, err := c.preparerForListRegionalBySubscriptionForTopicType(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalBySubscriptionForTopicType", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalBySubscriptionForTopicType", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListRegionalBySubscriptionForTopicType(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalBySubscriptionForTopicType", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListRegionalBySubscriptionForTopicTypeComplete retrieves all of the results into a single object
func (c EventSubscriptionsClient) ListRegionalBySubscriptionForTopicTypeComplete(ctx context.Context, id LocationTopicTypeId, options ListRegionalBySubscriptionForTopicTypeOptions) (ListRegionalBySubscriptionForTopicTypeCompleteResult, error) {
	return c.ListRegionalBySubscriptionForTopicTypeCompleteMatchingPredicate(ctx, id, options, EventSubscriptionPredicate{})
}

// ListRegionalBySubscriptionForTopicTypeCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EventSubscriptionsClient) ListRegionalBySubscriptionForTopicTypeCompleteMatchingPredicate(ctx context.Context, id LocationTopicTypeId, options ListRegionalBySubscriptionForTopicTypeOptions, predicate EventSubscriptionPredicate) (resp ListRegionalBySubscriptionForTopicTypeCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	page, err := c.ListRegionalBySubscriptionForTopicType(ctx, id, options)
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

	out := ListRegionalBySubscriptionForTopicTypeCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListRegionalBySubscriptionForTopicType prepares the ListRegionalBySubscriptionForTopicType request.
func (c EventSubscriptionsClient) preparerForListRegionalBySubscriptionForTopicType(ctx context.Context, id LocationTopicTypeId, options ListRegionalBySubscriptionForTopicTypeOptions) (*http.Request, error) {
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

// preparerForListRegionalBySubscriptionForTopicTypeWithNextLink prepares the ListRegionalBySubscriptionForTopicType request with the given nextLink token.
func (c EventSubscriptionsClient) preparerForListRegionalBySubscriptionForTopicTypeWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListRegionalBySubscriptionForTopicType handles the response to the ListRegionalBySubscriptionForTopicType request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForListRegionalBySubscriptionForTopicType(resp *http.Response) (result ListRegionalBySubscriptionForTopicTypeResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListRegionalBySubscriptionForTopicTypeResponse, err error) {
			req, err := c.preparerForListRegionalBySubscriptionForTopicTypeWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalBySubscriptionForTopicType", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalBySubscriptionForTopicType", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListRegionalBySubscriptionForTopicType(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalBySubscriptionForTopicType", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
