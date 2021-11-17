package eventchannels

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByPartnerNamespaceResponse struct {
	HttpResponse *http.Response
	Model        *[]EventChannel

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByPartnerNamespaceResponse, error)
}

type ListByPartnerNamespaceCompleteResult struct {
	Items []EventChannel
}

func (r ListByPartnerNamespaceResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByPartnerNamespaceResponse) LoadMore(ctx context.Context) (resp ListByPartnerNamespaceResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByPartnerNamespaceOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListByPartnerNamespaceOptions() ListByPartnerNamespaceOptions {
	return ListByPartnerNamespaceOptions{}
}

func (o ListByPartnerNamespaceOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListByPartnerNamespace ...
func (c EventChannelsClient) ListByPartnerNamespace(ctx context.Context, id PartnerNamespaceId, options ListByPartnerNamespaceOptions) (resp ListByPartnerNamespaceResponse, err error) {
	req, err := c.preparerForListByPartnerNamespace(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventchannels.EventChannelsClient", "ListByPartnerNamespace", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventchannels.EventChannelsClient", "ListByPartnerNamespace", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByPartnerNamespace(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventchannels.EventChannelsClient", "ListByPartnerNamespace", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByPartnerNamespaceComplete retrieves all of the results into a single object
func (c EventChannelsClient) ListByPartnerNamespaceComplete(ctx context.Context, id PartnerNamespaceId, options ListByPartnerNamespaceOptions) (ListByPartnerNamespaceCompleteResult, error) {
	return c.ListByPartnerNamespaceCompleteMatchingPredicate(ctx, id, options, EventChannelPredicate{})
}

// ListByPartnerNamespaceCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EventChannelsClient) ListByPartnerNamespaceCompleteMatchingPredicate(ctx context.Context, id PartnerNamespaceId, options ListByPartnerNamespaceOptions, predicate EventChannelPredicate) (resp ListByPartnerNamespaceCompleteResult, err error) {
	items := make([]EventChannel, 0)

	page, err := c.ListByPartnerNamespace(ctx, id, options)
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

	out := ListByPartnerNamespaceCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByPartnerNamespace prepares the ListByPartnerNamespace request.
func (c EventChannelsClient) preparerForListByPartnerNamespace(ctx context.Context, id PartnerNamespaceId, options ListByPartnerNamespaceOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/eventChannels", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByPartnerNamespaceWithNextLink prepares the ListByPartnerNamespace request with the given nextLink token.
func (c EventChannelsClient) preparerForListByPartnerNamespaceWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByPartnerNamespace handles the response to the ListByPartnerNamespace request. The method always
// closes the http.Response Body.
func (c EventChannelsClient) responderForListByPartnerNamespace(resp *http.Response) (result ListByPartnerNamespaceResponse, err error) {
	type page struct {
		Values   []EventChannel `json:"value"`
		NextLink *string        `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByPartnerNamespaceResponse, err error) {
			req, err := c.preparerForListByPartnerNamespaceWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventchannels.EventChannelsClient", "ListByPartnerNamespace", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventchannels.EventChannelsClient", "ListByPartnerNamespace", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByPartnerNamespace(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventchannels.EventChannelsClient", "ListByPartnerNamespace", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
