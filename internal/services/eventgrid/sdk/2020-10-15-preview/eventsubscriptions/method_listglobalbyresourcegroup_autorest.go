package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListGlobalByResourceGroupResponse struct {
	HttpResponse *http.Response
	Model        *[]EventSubscription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListGlobalByResourceGroupResponse, error)
}

type ListGlobalByResourceGroupCompleteResult struct {
	Items []EventSubscription
}

func (r ListGlobalByResourceGroupResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListGlobalByResourceGroupResponse) LoadMore(ctx context.Context) (resp ListGlobalByResourceGroupResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListGlobalByResourceGroupOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListGlobalByResourceGroupOptions() ListGlobalByResourceGroupOptions {
	return ListGlobalByResourceGroupOptions{}
}

func (o ListGlobalByResourceGroupOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListGlobalByResourceGroup ...
func (c EventSubscriptionsClient) ListGlobalByResourceGroup(ctx context.Context, id ResourceGroupId, options ListGlobalByResourceGroupOptions) (resp ListGlobalByResourceGroupResponse, err error) {
	req, err := c.preparerForListGlobalByResourceGroup(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListGlobalByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListGlobalByResourceGroupComplete retrieves all of the results into a single object
func (c EventSubscriptionsClient) ListGlobalByResourceGroupComplete(ctx context.Context, id ResourceGroupId, options ListGlobalByResourceGroupOptions) (ListGlobalByResourceGroupCompleteResult, error) {
	return c.ListGlobalByResourceGroupCompleteMatchingPredicate(ctx, id, options, EventSubscriptionPredicate{})
}

// ListGlobalByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EventSubscriptionsClient) ListGlobalByResourceGroupCompleteMatchingPredicate(ctx context.Context, id ResourceGroupId, options ListGlobalByResourceGroupOptions, predicate EventSubscriptionPredicate) (resp ListGlobalByResourceGroupCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	page, err := c.ListGlobalByResourceGroup(ctx, id, options)
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

	out := ListGlobalByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListGlobalByResourceGroup prepares the ListGlobalByResourceGroup request.
func (c EventSubscriptionsClient) preparerForListGlobalByResourceGroup(ctx context.Context, id ResourceGroupId, options ListGlobalByResourceGroupOptions) (*http.Request, error) {
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

// preparerForListGlobalByResourceGroupWithNextLink prepares the ListGlobalByResourceGroup request with the given nextLink token.
func (c EventSubscriptionsClient) preparerForListGlobalByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListGlobalByResourceGroup handles the response to the ListGlobalByResourceGroup request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForListGlobalByResourceGroup(resp *http.Response) (result ListGlobalByResourceGroupResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListGlobalByResourceGroupResponse, err error) {
			req, err := c.preparerForListGlobalByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListGlobalByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListGlobalByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
