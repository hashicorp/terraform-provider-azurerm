package confidentialledger

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type LedgerListByResourceGroupResponse struct {
	HttpResponse *http.Response
	Model        *[]ConfidentialLedger

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (LedgerListByResourceGroupResponse, error)
}

type LedgerListByResourceGroupCompleteResult struct {
	Items []ConfidentialLedger
}

func (r LedgerListByResourceGroupResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r LedgerListByResourceGroupResponse) LoadMore(ctx context.Context) (resp LedgerListByResourceGroupResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type LedgerListByResourceGroupOptions struct {
	Filter *string
}

func DefaultLedgerListByResourceGroupOptions() LedgerListByResourceGroupOptions {
	return LedgerListByResourceGroupOptions{}
}

func (o LedgerListByResourceGroupOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// LedgerListByResourceGroup ...
func (c ConfidentialLedgerClient) LedgerListByResourceGroup(ctx context.Context, id ResourceGroupId, options LedgerListByResourceGroupOptions) (resp LedgerListByResourceGroupResponse, err error) {
	req, err := c.preparerForLedgerListByResourceGroup(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForLedgerListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// LedgerListByResourceGroupComplete retrieves all of the results into a single object
func (c ConfidentialLedgerClient) LedgerListByResourceGroupComplete(ctx context.Context, id ResourceGroupId, options LedgerListByResourceGroupOptions) (LedgerListByResourceGroupCompleteResult, error) {
	return c.LedgerListByResourceGroupCompleteMatchingPredicate(ctx, id, options, ConfidentialLedgerPredicate{})
}

// LedgerListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ConfidentialLedgerClient) LedgerListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id ResourceGroupId, options LedgerListByResourceGroupOptions, predicate ConfidentialLedgerPredicate) (resp LedgerListByResourceGroupCompleteResult, err error) {
	items := make([]ConfidentialLedger, 0)

	page, err := c.LedgerListByResourceGroup(ctx, id, options)
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

	out := LedgerListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForLedgerListByResourceGroup prepares the LedgerListByResourceGroup request.
func (c ConfidentialLedgerClient) preparerForLedgerListByResourceGroup(ctx context.Context, id ResourceGroupId, options LedgerListByResourceGroupOptions) (*http.Request, error) {
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

// preparerForLedgerListByResourceGroupWithNextLink prepares the LedgerListByResourceGroup request with the given nextLink token.
func (c ConfidentialLedgerClient) preparerForLedgerListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForLedgerListByResourceGroup handles the response to the LedgerListByResourceGroup request. The method always
// closes the http.Response Body.
func (c ConfidentialLedgerClient) responderForLedgerListByResourceGroup(resp *http.Response) (result LedgerListByResourceGroupResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result LedgerListByResourceGroupResponse, err error) {
			req, err := c.preparerForLedgerListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForLedgerListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "confidentialledger.ConfidentialLedgerClient", "LedgerListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
