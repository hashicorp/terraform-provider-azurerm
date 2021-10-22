package get

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PrivateEndpointConnectionsListResponse struct {
	HttpResponse *http.Response
	Model        *[]PrivateEndpointConnection

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (PrivateEndpointConnectionsListResponse, error)
}

type PrivateEndpointConnectionsListCompleteResult struct {
	Items []PrivateEndpointConnection
}

func (r PrivateEndpointConnectionsListResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r PrivateEndpointConnectionsListResponse) LoadMore(ctx context.Context) (resp PrivateEndpointConnectionsListResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// PrivateEndpointConnectionsList ...
func (c GETClient) PrivateEndpointConnectionsList(ctx context.Context, id WorkspaceId) (resp PrivateEndpointConnectionsListResponse, err error) {
	req, err := c.preparerForPrivateEndpointConnectionsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "get.GETClient", "PrivateEndpointConnectionsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "get.GETClient", "PrivateEndpointConnectionsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForPrivateEndpointConnectionsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "get.GETClient", "PrivateEndpointConnectionsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// PrivateEndpointConnectionsListComplete retrieves all of the results into a single object
func (c GETClient) PrivateEndpointConnectionsListComplete(ctx context.Context, id WorkspaceId) (PrivateEndpointConnectionsListCompleteResult, error) {
	return c.PrivateEndpointConnectionsListCompleteMatchingPredicate(ctx, id, PrivateEndpointConnectionPredicate{})
}

// PrivateEndpointConnectionsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c GETClient) PrivateEndpointConnectionsListCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, predicate PrivateEndpointConnectionPredicate) (resp PrivateEndpointConnectionsListCompleteResult, err error) {
	items := make([]PrivateEndpointConnection, 0)

	page, err := c.PrivateEndpointConnectionsList(ctx, id)
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

	out := PrivateEndpointConnectionsListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForPrivateEndpointConnectionsList prepares the PrivateEndpointConnectionsList request.
func (c GETClient) preparerForPrivateEndpointConnectionsList(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/privateEndpointConnections", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForPrivateEndpointConnectionsListWithNextLink prepares the PrivateEndpointConnectionsList request with the given nextLink token.
func (c GETClient) preparerForPrivateEndpointConnectionsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForPrivateEndpointConnectionsList handles the response to the PrivateEndpointConnectionsList request. The method always
// closes the http.Response Body.
func (c GETClient) responderForPrivateEndpointConnectionsList(resp *http.Response) (result PrivateEndpointConnectionsListResponse, err error) {
	type page struct {
		Values   []PrivateEndpointConnection `json:"value"`
		NextLink *string                     `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result PrivateEndpointConnectionsListResponse, err error) {
			req, err := c.preparerForPrivateEndpointConnectionsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "get.GETClient", "PrivateEndpointConnectionsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "get.GETClient", "PrivateEndpointConnectionsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForPrivateEndpointConnectionsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "get.GETClient", "PrivateEndpointConnectionsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
