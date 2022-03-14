package vmhhostlist

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type VMHostListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]VMResources

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (VMHostListOperationResponse, error)
}

type VMHostListCompleteResult struct {
	Items []VMResources
}

func (r VMHostListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r VMHostListOperationResponse) LoadMore(ctx context.Context) (resp VMHostListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// VMHostList ...
func (c VMHHostListClient) VMHostList(ctx context.Context, id MonitorId) (resp VMHostListOperationResponse, err error) {
	req, err := c.preparerForVMHostList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vmhhostlist.VMHHostListClient", "VMHostList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "vmhhostlist.VMHHostListClient", "VMHostList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForVMHostList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vmhhostlist.VMHHostListClient", "VMHostList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// VMHostListComplete retrieves all of the results into a single object
func (c VMHHostListClient) VMHostListComplete(ctx context.Context, id MonitorId) (VMHostListCompleteResult, error) {
	return c.VMHostListCompleteMatchingPredicate(ctx, id, VMResourcesOperationPredicate{})
}

// VMHostListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c VMHHostListClient) VMHostListCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate VMResourcesOperationPredicate) (resp VMHostListCompleteResult, err error) {
	items := make([]VMResources, 0)

	page, err := c.VMHostList(ctx, id)
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

	out := VMHostListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForVMHostList prepares the VMHostList request.
func (c VMHHostListClient) preparerForVMHostList(ctx context.Context, id MonitorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listVMHost", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForVMHostListWithNextLink prepares the VMHostList request with the given nextLink token.
func (c VMHHostListClient) preparerForVMHostListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForVMHostList handles the response to the VMHostList request. The method always
// closes the http.Response Body.
func (c VMHHostListClient) responderForVMHostList(resp *http.Response) (result VMHostListOperationResponse, err error) {
	type page struct {
		Values   []VMResources `json:"value"`
		NextLink *string       `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result VMHostListOperationResponse, err error) {
			req, err := c.preparerForVMHostListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "vmhhostlist.VMHHostListClient", "VMHostList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "vmhhostlist.VMHHostListClient", "VMHostList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForVMHostList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "vmhhostlist.VMHHostListClient", "VMHostList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
