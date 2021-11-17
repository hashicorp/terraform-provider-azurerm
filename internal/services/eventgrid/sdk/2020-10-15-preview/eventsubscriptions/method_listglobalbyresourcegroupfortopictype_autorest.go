package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListGlobalByResourceGroupForTopicTypeResponse struct {
	HttpResponse *http.Response
	Model        *[]EventSubscription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListGlobalByResourceGroupForTopicTypeResponse, error)
}

type ListGlobalByResourceGroupForTopicTypeCompleteResult struct {
	Items []EventSubscription
}

func (r ListGlobalByResourceGroupForTopicTypeResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListGlobalByResourceGroupForTopicTypeResponse) LoadMore(ctx context.Context) (resp ListGlobalByResourceGroupForTopicTypeResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListGlobalByResourceGroupForTopicTypeOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListGlobalByResourceGroupForTopicTypeOptions() ListGlobalByResourceGroupForTopicTypeOptions {
	return ListGlobalByResourceGroupForTopicTypeOptions{}
}

func (o ListGlobalByResourceGroupForTopicTypeOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListGlobalByResourceGroupForTopicType ...
func (c EventSubscriptionsClient) ListGlobalByResourceGroupForTopicType(ctx context.Context, id ResourceGroupProviderTopicTypeId, options ListGlobalByResourceGroupForTopicTypeOptions) (resp ListGlobalByResourceGroupForTopicTypeResponse, err error) {
	req, err := c.preparerForListGlobalByResourceGroupForTopicType(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalByResourceGroupForTopicType", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalByResourceGroupForTopicType", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListGlobalByResourceGroupForTopicType(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalByResourceGroupForTopicType", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListGlobalByResourceGroupForTopicTypeComplete retrieves all of the results into a single object
func (c EventSubscriptionsClient) ListGlobalByResourceGroupForTopicTypeComplete(ctx context.Context, id ResourceGroupProviderTopicTypeId, options ListGlobalByResourceGroupForTopicTypeOptions) (ListGlobalByResourceGroupForTopicTypeCompleteResult, error) {
	return c.ListGlobalByResourceGroupForTopicTypeCompleteMatchingPredicate(ctx, id, options, EventSubscriptionPredicate{})
}

// ListGlobalByResourceGroupForTopicTypeCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EventSubscriptionsClient) ListGlobalByResourceGroupForTopicTypeCompleteMatchingPredicate(ctx context.Context, id ResourceGroupProviderTopicTypeId, options ListGlobalByResourceGroupForTopicTypeOptions, predicate EventSubscriptionPredicate) (resp ListGlobalByResourceGroupForTopicTypeCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	page, err := c.ListGlobalByResourceGroupForTopicType(ctx, id, options)
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

	out := ListGlobalByResourceGroupForTopicTypeCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListGlobalByResourceGroupForTopicType prepares the ListGlobalByResourceGroupForTopicType request.
func (c EventSubscriptionsClient) preparerForListGlobalByResourceGroupForTopicType(ctx context.Context, id ResourceGroupProviderTopicTypeId, options ListGlobalByResourceGroupForTopicTypeOptions) (*http.Request, error) {
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

// preparerForListGlobalByResourceGroupForTopicTypeWithNextLink prepares the ListGlobalByResourceGroupForTopicType request with the given nextLink token.
func (c EventSubscriptionsClient) preparerForListGlobalByResourceGroupForTopicTypeWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListGlobalByResourceGroupForTopicType handles the response to the ListGlobalByResourceGroupForTopicType request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForListGlobalByResourceGroupForTopicType(resp *http.Response) (result ListGlobalByResourceGroupForTopicTypeResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListGlobalByResourceGroupForTopicTypeResponse, err error) {
			req, err := c.preparerForListGlobalByResourceGroupForTopicTypeWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalByResourceGroupForTopicType", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalByResourceGroupForTopicType", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListGlobalByResourceGroupForTopicType(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalByResourceGroupForTopicType", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
