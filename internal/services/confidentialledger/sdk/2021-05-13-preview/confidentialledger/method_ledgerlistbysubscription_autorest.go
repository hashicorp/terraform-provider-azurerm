package confidentialledger

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type LedgerListBySubscriptionResponse struct {
	HttpResponse *http.Response
	Model        *[]ConfidentialLedger

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (LedgerListBySubscriptionResponse, error)
}

type LedgerListBySubscriptionCompleteResult struct {
	Items []ConfidentialLedger
}

func (r LedgerListBySubscriptionResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r LedgerListBySubscriptionResponse) LoadMore(ctx context.Context) (resp LedgerListBySubscriptionResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type LedgerListBySubscriptionOptions struct {
	Filter *string
}

func DefaultLedgerListBySubscriptionOptions() LedgerListBySubscriptionOptions {
	return LedgerListBySubscriptionOptions{}
}

func (o LedgerListBySubscriptionOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// LedgerListBySubscription ...
func (c ConfidentialLedgerClient) LedgerListBySubscription(ctx context.Context, id SubscriptionId, options LedgerListBySubscriptionOptions) (resp LedgerListBySubscriptionResponse, err error) {
	req, err := c.preparerForLedgerListBySubscription(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerListBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerListBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForLedgerListBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerListBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// LedgerListBySubscriptionComplete retrieves all of the results into a single object
func (c ConfidentialLedgerClient) LedgerListBySubscriptionComplete(ctx context.Context, id SubscriptionId, options LedgerListBySubscriptionOptions) (LedgerListBySubscriptionCompleteResult, error) {
	return c.LedgerListBySubscriptionCompleteMatchingPredicate(ctx, id, options, ConfidentialLedgerPredicate{})
}

// LedgerListBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ConfidentialLedgerClient) LedgerListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id SubscriptionId, options LedgerListBySubscriptionOptions, predicate ConfidentialLedgerPredicate) (resp LedgerListBySubscriptionCompleteResult, err error) {
	items := make([]ConfidentialLedger, 0)

	page, err := c.LedgerListBySubscription(ctx, id, options)
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

	out := LedgerListBySubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForLedgerListBySubscription prepares the LedgerListBySubscription request.
func (c ConfidentialLedgerClient) preparerForLedgerListBySubscription(ctx context.Context, id SubscriptionId, options LedgerListBySubscriptionOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.ConfidentialLedger/ledgers", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForLedgerListBySubscriptionWithNextLink prepares the LedgerListBySubscription request with the given nextLink token.
func (c ConfidentialLedgerClient) preparerForLedgerListBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForLedgerListBySubscription handles the response to the LedgerListBySubscription request. The method always
// closes the http.Response Body.
func (c ConfidentialLedgerClient) responderForLedgerListBySubscription(resp *http.Response) (result LedgerListBySubscriptionResponse, err error) {
	type page struct {
		Values   []ConfidentialLedger `json:"value"`
		NextLink *string              `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result LedgerListBySubscriptionResponse, err error) {
			req, err := c.preparerForLedgerListBySubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerListBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerListBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForLedgerListBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerListBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
