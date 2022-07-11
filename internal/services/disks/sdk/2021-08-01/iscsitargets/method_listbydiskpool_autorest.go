package iscsitargets

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByDiskPoolResponse struct {
	HttpResponse *http.Response
	Model        *[]IscsiTarget

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByDiskPoolResponse, error)
}

type ListByDiskPoolCompleteResult struct {
	Items []IscsiTarget
}

func (r ListByDiskPoolResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByDiskPoolResponse) LoadMore(ctx context.Context) (resp ListByDiskPoolResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByDiskPool ...
func (c IscsiTargetsClient) ListByDiskPool(ctx context.Context, id DiskPoolId) (resp ListByDiskPoolResponse, err error) {
	req, err := c.preparerForListByDiskPool(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iscsitargets.IscsiTargetsClient", "ListByDiskPool", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "iscsitargets.IscsiTargetsClient", "ListByDiskPool", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByDiskPool(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iscsitargets.IscsiTargetsClient", "ListByDiskPool", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByDiskPoolComplete retrieves all of the results into a single object
func (c IscsiTargetsClient) ListByDiskPoolComplete(ctx context.Context, id DiskPoolId) (ListByDiskPoolCompleteResult, error) {
	return c.ListByDiskPoolCompleteMatchingPredicate(ctx, id, IscsiTargetPredicate{})
}

// ListByDiskPoolCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c IscsiTargetsClient) ListByDiskPoolCompleteMatchingPredicate(ctx context.Context, id DiskPoolId, predicate IscsiTargetPredicate) (resp ListByDiskPoolCompleteResult, err error) {
	items := make([]IscsiTarget, 0)

	page, err := c.ListByDiskPool(ctx, id)
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

	out := ListByDiskPoolCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByDiskPool prepares the ListByDiskPool request.
func (c IscsiTargetsClient) preparerForListByDiskPool(ctx context.Context, id DiskPoolId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/iscsiTargets", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByDiskPoolWithNextLink prepares the ListByDiskPool request with the given nextLink token.
func (c IscsiTargetsClient) preparerForListByDiskPoolWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByDiskPool handles the response to the ListByDiskPool request. The method always
// closes the http.Response Body.
func (c IscsiTargetsClient) responderForListByDiskPool(resp *http.Response) (result ListByDiskPoolResponse, err error) {
	type page struct {
		Values   []IscsiTarget `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByDiskPoolResponse, err error) {
			req, err := c.preparerForListByDiskPoolWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "iscsitargets.IscsiTargetsClient", "ListByDiskPool", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "iscsitargets.IscsiTargetsClient", "ListByDiskPool", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByDiskPool(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "iscsitargets.IscsiTargetsClient", "ListByDiskPool", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
