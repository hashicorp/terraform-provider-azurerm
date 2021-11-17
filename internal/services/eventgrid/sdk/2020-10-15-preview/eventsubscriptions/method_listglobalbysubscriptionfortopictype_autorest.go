package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListGlobalBySubscriptionForTopicTypeResponse struct {
	HttpResponse *http.Response
	Model        *[]EventSubscription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListGlobalBySubscriptionForTopicTypeResponse, error)
}

type ListGlobalBySubscriptionForTopicTypeCompleteResult struct {
	Items []EventSubscription
}

func (r ListGlobalBySubscriptionForTopicTypeResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListGlobalBySubscriptionForTopicTypeResponse) LoadMore(ctx context.Context) (resp ListGlobalBySubscriptionForTopicTypeResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListGlobalBySubscriptionForTopicTypeOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListGlobalBySubscriptionForTopicTypeOptions() ListGlobalBySubscriptionForTopicTypeOptions {
	return ListGlobalBySubscriptionForTopicTypeOptions{}
}

func (o ListGlobalBySubscriptionForTopicTypeOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListGlobalBySubscriptionForTopicType ...
func (c EventSubscriptionsClient) ListGlobalBySubscriptionForTopicType(ctx context.Context, id ProviderTopicTypeId, options ListGlobalBySubscriptionForTopicTypeOptions) (resp ListGlobalBySubscriptionForTopicTypeResponse, err error) {
	req, err := c.preparerForListGlobalBySubscriptionForTopicType(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalBySubscriptionForTopicType", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalBySubscriptionForTopicType", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListGlobalBySubscriptionForTopicType(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalBySubscriptionForTopicType", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListGlobalBySubscriptionForTopicTypeComplete retrieves all of the results into a single object
func (c EventSubscriptionsClient) ListGlobalBySubscriptionForTopicTypeComplete(ctx context.Context, id ProviderTopicTypeId, options ListGlobalBySubscriptionForTopicTypeOptions) (ListGlobalBySubscriptionForTopicTypeCompleteResult, error) {
	return c.ListGlobalBySubscriptionForTopicTypeCompleteMatchingPredicate(ctx, id, options, EventSubscriptionPredicate{})
}

// ListGlobalBySubscriptionForTopicTypeCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EventSubscriptionsClient) ListGlobalBySubscriptionForTopicTypeCompleteMatchingPredicate(ctx context.Context, id ProviderTopicTypeId, options ListGlobalBySubscriptionForTopicTypeOptions, predicate EventSubscriptionPredicate) (resp ListGlobalBySubscriptionForTopicTypeCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	page, err := c.ListGlobalBySubscriptionForTopicType(ctx, id, options)
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

	out := ListGlobalBySubscriptionForTopicTypeCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListGlobalBySubscriptionForTopicType prepares the ListGlobalBySubscriptionForTopicType request.
func (c EventSubscriptionsClient) preparerForListGlobalBySubscriptionForTopicType(ctx context.Context, id ProviderTopicTypeId, options ListGlobalBySubscriptionForTopicTypeOptions) (*http.Request, error) {
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

// preparerForListGlobalBySubscriptionForTopicTypeWithNextLink prepares the ListGlobalBySubscriptionForTopicType request with the given nextLink token.
func (c EventSubscriptionsClient) preparerForListGlobalBySubscriptionForTopicTypeWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListGlobalBySubscriptionForTopicType handles the response to the ListGlobalBySubscriptionForTopicType request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForListGlobalBySubscriptionForTopicType(resp *http.Response) (result ListGlobalBySubscriptionForTopicTypeResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListGlobalBySubscriptionForTopicTypeResponse, err error) {
			req, err := c.preparerForListGlobalBySubscriptionForTopicTypeWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalBySubscriptionForTopicType", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalBySubscriptionForTopicType", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListGlobalBySubscriptionForTopicType(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalBySubscriptionForTopicType", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
