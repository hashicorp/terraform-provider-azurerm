package vnetpeering

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByWorkspaceResponse struct {
	HttpResponse *http.Response
	Model        *[]VirtualNetworkPeering

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByWorkspaceResponse, error)
}

type ListByWorkspaceCompleteResult struct {
	Items []VirtualNetworkPeering
}

func (r ListByWorkspaceResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByWorkspaceResponse) LoadMore(ctx context.Context) (resp ListByWorkspaceResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByWorkspace ...
func (c VNetPeeringClient) ListByWorkspace(ctx context.Context, id WorkspaceId) (resp ListByWorkspaceResponse, err error) {
	req, err := c.preparerForListByWorkspace(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vnetpeering.VNetPeeringClient", "ListByWorkspace", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "vnetpeering.VNetPeeringClient", "ListByWorkspace", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByWorkspace(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vnetpeering.VNetPeeringClient", "ListByWorkspace", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByWorkspaceComplete retrieves all of the results into a single object
func (c VNetPeeringClient) ListByWorkspaceComplete(ctx context.Context, id WorkspaceId) (ListByWorkspaceCompleteResult, error) {
	return c.ListByWorkspaceCompleteMatchingPredicate(ctx, id, VirtualNetworkPeeringPredicate{})
}

// ListByWorkspaceCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c VNetPeeringClient) ListByWorkspaceCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, predicate VirtualNetworkPeeringPredicate) (resp ListByWorkspaceCompleteResult, err error) {
	items := make([]VirtualNetworkPeering, 0)

	page, err := c.ListByWorkspace(ctx, id)
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

	out := ListByWorkspaceCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByWorkspace prepares the ListByWorkspace request.
func (c VNetPeeringClient) preparerForListByWorkspace(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/virtualNetworkPeerings", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByWorkspaceWithNextLink prepares the ListByWorkspace request with the given nextLink token.
func (c VNetPeeringClient) preparerForListByWorkspaceWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByWorkspace handles the response to the ListByWorkspace request. The method always
// closes the http.Response Body.
func (c VNetPeeringClient) responderForListByWorkspace(resp *http.Response) (result ListByWorkspaceResponse, err error) {
	type page struct {
		Values   []VirtualNetworkPeering `json:"value"`
		NextLink *string                 `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByWorkspaceResponse, err error) {
			req, err := c.preparerForListByWorkspaceWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "vnetpeering.VNetPeeringClient", "ListByWorkspace", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "vnetpeering.VNetPeeringClient", "ListByWorkspace", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByWorkspace(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "vnetpeering.VNetPeeringClient", "ListByWorkspace", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
