package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type SystemTopicEventSubscriptionsListBySystemTopicResponse struct {
	HttpResponse *http.Response
	Model        *[]EventSubscription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (SystemTopicEventSubscriptionsListBySystemTopicResponse, error)
}

type SystemTopicEventSubscriptionsListBySystemTopicCompleteResult struct {
	Items []EventSubscription
}

func (r SystemTopicEventSubscriptionsListBySystemTopicResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r SystemTopicEventSubscriptionsListBySystemTopicResponse) LoadMore(ctx context.Context) (resp SystemTopicEventSubscriptionsListBySystemTopicResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type SystemTopicEventSubscriptionsListBySystemTopicOptions struct {
	Filter *string
	Top    *int64
}

func DefaultSystemTopicEventSubscriptionsListBySystemTopicOptions() SystemTopicEventSubscriptionsListBySystemTopicOptions {
	return SystemTopicEventSubscriptionsListBySystemTopicOptions{}
}

func (o SystemTopicEventSubscriptionsListBySystemTopicOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// SystemTopicEventSubscriptionsListBySystemTopic ...
func (c EventSubscriptionsClient) SystemTopicEventSubscriptionsListBySystemTopic(ctx context.Context, id SystemTopicId, options SystemTopicEventSubscriptionsListBySystemTopicOptions) (resp SystemTopicEventSubscriptionsListBySystemTopicResponse, err error) {
	req, err := c.preparerForSystemTopicEventSubscriptionsListBySystemTopic(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsListBySystemTopic", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsListBySystemTopic", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForSystemTopicEventSubscriptionsListBySystemTopic(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsListBySystemTopic", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// SystemTopicEventSubscriptionsListBySystemTopicComplete retrieves all of the results into a single object
func (c EventSubscriptionsClient) SystemTopicEventSubscriptionsListBySystemTopicComplete(ctx context.Context, id SystemTopicId, options SystemTopicEventSubscriptionsListBySystemTopicOptions) (SystemTopicEventSubscriptionsListBySystemTopicCompleteResult, error) {
	return c.SystemTopicEventSubscriptionsListBySystemTopicCompleteMatchingPredicate(ctx, id, options, EventSubscriptionPredicate{})
}

// SystemTopicEventSubscriptionsListBySystemTopicCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EventSubscriptionsClient) SystemTopicEventSubscriptionsListBySystemTopicCompleteMatchingPredicate(ctx context.Context, id SystemTopicId, options SystemTopicEventSubscriptionsListBySystemTopicOptions, predicate EventSubscriptionPredicate) (resp SystemTopicEventSubscriptionsListBySystemTopicCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	page, err := c.SystemTopicEventSubscriptionsListBySystemTopic(ctx, id, options)
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

	out := SystemTopicEventSubscriptionsListBySystemTopicCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForSystemTopicEventSubscriptionsListBySystemTopic prepares the SystemTopicEventSubscriptionsListBySystemTopic request.
func (c EventSubscriptionsClient) preparerForSystemTopicEventSubscriptionsListBySystemTopic(ctx context.Context, id SystemTopicId, options SystemTopicEventSubscriptionsListBySystemTopicOptions) (*http.Request, error) {
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

// preparerForSystemTopicEventSubscriptionsListBySystemTopicWithNextLink prepares the SystemTopicEventSubscriptionsListBySystemTopic request with the given nextLink token.
func (c EventSubscriptionsClient) preparerForSystemTopicEventSubscriptionsListBySystemTopicWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForSystemTopicEventSubscriptionsListBySystemTopic handles the response to the SystemTopicEventSubscriptionsListBySystemTopic request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForSystemTopicEventSubscriptionsListBySystemTopic(resp *http.Response) (result SystemTopicEventSubscriptionsListBySystemTopicResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result SystemTopicEventSubscriptionsListBySystemTopicResponse, err error) {
			req, err := c.preparerForSystemTopicEventSubscriptionsListBySystemTopicWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsListBySystemTopic", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsListBySystemTopic", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForSystemTopicEventSubscriptionsListBySystemTopic(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "SystemTopicEventSubscriptionsListBySystemTopic", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
