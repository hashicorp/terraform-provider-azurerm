package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListRegionalByResourceGroupForTopicTypeResponse struct {
	HttpResponse *http.Response
	Model        *[]EventSubscription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListRegionalByResourceGroupForTopicTypeResponse, error)
}

type ListRegionalByResourceGroupForTopicTypeCompleteResult struct {
	Items []EventSubscription
}

func (r ListRegionalByResourceGroupForTopicTypeResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListRegionalByResourceGroupForTopicTypeResponse) LoadMore(ctx context.Context) (resp ListRegionalByResourceGroupForTopicTypeResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListRegionalByResourceGroupForTopicTypeOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListRegionalByResourceGroupForTopicTypeOptions() ListRegionalByResourceGroupForTopicTypeOptions {
	return ListRegionalByResourceGroupForTopicTypeOptions{}
}

func (o ListRegionalByResourceGroupForTopicTypeOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListRegionalByResourceGroupForTopicType ...
func (c EventSubscriptionsClient) ListRegionalByResourceGroupForTopicType(ctx context.Context, id ProviderLocationTopicTypeId, options ListRegionalByResourceGroupForTopicTypeOptions) (resp ListRegionalByResourceGroupForTopicTypeResponse, err error) {
	req, err := c.preparerForListRegionalByResourceGroupForTopicType(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalByResourceGroupForTopicType", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalByResourceGroupForTopicType", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListRegionalByResourceGroupForTopicType(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalByResourceGroupForTopicType", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListRegionalByResourceGroupForTopicTypeComplete retrieves all of the results into a single object
func (c EventSubscriptionsClient) ListRegionalByResourceGroupForTopicTypeComplete(ctx context.Context, id ProviderLocationTopicTypeId, options ListRegionalByResourceGroupForTopicTypeOptions) (ListRegionalByResourceGroupForTopicTypeCompleteResult, error) {
	return c.ListRegionalByResourceGroupForTopicTypeCompleteMatchingPredicate(ctx, id, options, EventSubscriptionPredicate{})
}

// ListRegionalByResourceGroupForTopicTypeCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EventSubscriptionsClient) ListRegionalByResourceGroupForTopicTypeCompleteMatchingPredicate(ctx context.Context, id ProviderLocationTopicTypeId, options ListRegionalByResourceGroupForTopicTypeOptions, predicate EventSubscriptionPredicate) (resp ListRegionalByResourceGroupForTopicTypeCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	page, err := c.ListRegionalByResourceGroupForTopicType(ctx, id, options)
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

	out := ListRegionalByResourceGroupForTopicTypeCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListRegionalByResourceGroupForTopicType prepares the ListRegionalByResourceGroupForTopicType request.
func (c EventSubscriptionsClient) preparerForListRegionalByResourceGroupForTopicType(ctx context.Context, id ProviderLocationTopicTypeId, options ListRegionalByResourceGroupForTopicTypeOptions) (*http.Request, error) {
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

// preparerForListRegionalByResourceGroupForTopicTypeWithNextLink prepares the ListRegionalByResourceGroupForTopicType request with the given nextLink token.
func (c EventSubscriptionsClient) preparerForListRegionalByResourceGroupForTopicTypeWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListRegionalByResourceGroupForTopicType handles the response to the ListRegionalByResourceGroupForTopicType request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForListRegionalByResourceGroupForTopicType(resp *http.Response) (result ListRegionalByResourceGroupForTopicTypeResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListRegionalByResourceGroupForTopicTypeResponse, err error) {
			req, err := c.preparerForListRegionalByResourceGroupForTopicTypeWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalByResourceGroupForTopicType", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalByResourceGroupForTopicType", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListRegionalByResourceGroupForTopicType(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListRegionalByResourceGroupForTopicType", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
