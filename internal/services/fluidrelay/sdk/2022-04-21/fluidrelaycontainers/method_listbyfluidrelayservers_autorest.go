package fluidrelaycontainers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByFluidRelayServersOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]FluidRelayContainer

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByFluidRelayServersOperationResponse, error)
}

type ListByFluidRelayServersCompleteResult struct {
	Items []FluidRelayContainer
}

func (r ListByFluidRelayServersOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByFluidRelayServersOperationResponse) LoadMore(ctx context.Context) (resp ListByFluidRelayServersOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByFluidRelayServers ...
func (c FluidRelayContainersClient) ListByFluidRelayServers(ctx context.Context, id FluidRelayServerId) (resp ListByFluidRelayServersOperationResponse, err error) {
	req, err := c.preparerForListByFluidRelayServers(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "fluidrelaycontainers.FluidRelayContainersClient", "ListByFluidRelayServers", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "fluidrelaycontainers.FluidRelayContainersClient", "ListByFluidRelayServers", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByFluidRelayServers(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "fluidrelaycontainers.FluidRelayContainersClient", "ListByFluidRelayServers", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByFluidRelayServersComplete retrieves all of the results into a single object
func (c FluidRelayContainersClient) ListByFluidRelayServersComplete(ctx context.Context, id FluidRelayServerId) (ListByFluidRelayServersCompleteResult, error) {
	return c.ListByFluidRelayServersCompleteMatchingPredicate(ctx, id, FluidRelayContainerOperationPredicate{})
}

// ListByFluidRelayServersCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c FluidRelayContainersClient) ListByFluidRelayServersCompleteMatchingPredicate(ctx context.Context, id FluidRelayServerId, predicate FluidRelayContainerOperationPredicate) (resp ListByFluidRelayServersCompleteResult, err error) {
	items := make([]FluidRelayContainer, 0)

	page, err := c.ListByFluidRelayServers(ctx, id)
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

	out := ListByFluidRelayServersCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByFluidRelayServers prepares the ListByFluidRelayServers request.
func (c FluidRelayContainersClient) preparerForListByFluidRelayServers(ctx context.Context, id FluidRelayServerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/fluidRelayContainers", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByFluidRelayServersWithNextLink prepares the ListByFluidRelayServers request with the given nextLink token.
func (c FluidRelayContainersClient) preparerForListByFluidRelayServersWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByFluidRelayServers handles the response to the ListByFluidRelayServers request. The method always
// closes the http.Response Body.
func (c FluidRelayContainersClient) responderForListByFluidRelayServers(resp *http.Response) (result ListByFluidRelayServersOperationResponse, err error) {
	type page struct {
		Values   []FluidRelayContainer `json:"value"`
		NextLink *string               `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByFluidRelayServersOperationResponse, err error) {
			req, err := c.preparerForListByFluidRelayServersWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "fluidrelaycontainers.FluidRelayContainersClient", "ListByFluidRelayServers", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "fluidrelaycontainers.FluidRelayContainersClient", "ListByFluidRelayServers", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByFluidRelayServers(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "fluidrelaycontainers.FluidRelayContainersClient", "ListByFluidRelayServers", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
