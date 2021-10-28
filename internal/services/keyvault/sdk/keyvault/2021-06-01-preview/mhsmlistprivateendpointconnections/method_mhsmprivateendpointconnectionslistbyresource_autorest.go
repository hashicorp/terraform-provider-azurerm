package mhsmlistprivateendpointconnections

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type MHSMPrivateEndpointConnectionsListByResourceResponse struct {
	HttpResponse *http.Response
	Model        *[]MHSMPrivateEndpointConnection

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (MHSMPrivateEndpointConnectionsListByResourceResponse, error)
}

type MHSMPrivateEndpointConnectionsListByResourceCompleteResult struct {
	Items []MHSMPrivateEndpointConnection
}

func (r MHSMPrivateEndpointConnectionsListByResourceResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r MHSMPrivateEndpointConnectionsListByResourceResponse) LoadMore(ctx context.Context) (resp MHSMPrivateEndpointConnectionsListByResourceResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// MHSMPrivateEndpointConnectionsListByResource ...
func (c MHSMListPrivateEndpointConnectionsClient) MHSMPrivateEndpointConnectionsListByResource(ctx context.Context, id ManagedHSMId) (resp MHSMPrivateEndpointConnectionsListByResourceResponse, err error) {
	req, err := c.preparerForMHSMPrivateEndpointConnectionsListByResource(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mhsmlistprivateendpointconnections.MHSMListPrivateEndpointConnectionsClient", "MHSMPrivateEndpointConnectionsListByResource", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "mhsmlistprivateendpointconnections.MHSMListPrivateEndpointConnectionsClient", "MHSMPrivateEndpointConnectionsListByResource", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForMHSMPrivateEndpointConnectionsListByResource(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mhsmlistprivateendpointconnections.MHSMListPrivateEndpointConnectionsClient", "MHSMPrivateEndpointConnectionsListByResource", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// MHSMPrivateEndpointConnectionsListByResourceComplete retrieves all of the results into a single object
func (c MHSMListPrivateEndpointConnectionsClient) MHSMPrivateEndpointConnectionsListByResourceComplete(ctx context.Context, id ManagedHSMId) (MHSMPrivateEndpointConnectionsListByResourceCompleteResult, error) {
	return c.MHSMPrivateEndpointConnectionsListByResourceCompleteMatchingPredicate(ctx, id, MHSMPrivateEndpointConnectionPredicate{})
}

// MHSMPrivateEndpointConnectionsListByResourceCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c MHSMListPrivateEndpointConnectionsClient) MHSMPrivateEndpointConnectionsListByResourceCompleteMatchingPredicate(ctx context.Context, id ManagedHSMId, predicate MHSMPrivateEndpointConnectionPredicate) (resp MHSMPrivateEndpointConnectionsListByResourceCompleteResult, err error) {
	items := make([]MHSMPrivateEndpointConnection, 0)

	page, err := c.MHSMPrivateEndpointConnectionsListByResource(ctx, id)
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

	out := MHSMPrivateEndpointConnectionsListByResourceCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForMHSMPrivateEndpointConnectionsListByResource prepares the MHSMPrivateEndpointConnectionsListByResource request.
func (c MHSMListPrivateEndpointConnectionsClient) preparerForMHSMPrivateEndpointConnectionsListByResource(ctx context.Context, id ManagedHSMId) (*http.Request, error) {
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

// preparerForMHSMPrivateEndpointConnectionsListByResourceWithNextLink prepares the MHSMPrivateEndpointConnectionsListByResource request with the given nextLink token.
func (c MHSMListPrivateEndpointConnectionsClient) preparerForMHSMPrivateEndpointConnectionsListByResourceWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForMHSMPrivateEndpointConnectionsListByResource handles the response to the MHSMPrivateEndpointConnectionsListByResource request. The method always
// closes the http.Response Body.
func (c MHSMListPrivateEndpointConnectionsClient) responderForMHSMPrivateEndpointConnectionsListByResource(resp *http.Response) (result MHSMPrivateEndpointConnectionsListByResourceResponse, err error) {
	type page struct {
		Values   []MHSMPrivateEndpointConnection `json:"value"`
		NextLink *string                         `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result MHSMPrivateEndpointConnectionsListByResourceResponse, err error) {
			req, err := c.preparerForMHSMPrivateEndpointConnectionsListByResourceWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "mhsmlistprivateendpointconnections.MHSMListPrivateEndpointConnectionsClient", "MHSMPrivateEndpointConnectionsListByResource", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "mhsmlistprivateendpointconnections.MHSMListPrivateEndpointConnectionsClient", "MHSMPrivateEndpointConnectionsListByResource", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForMHSMPrivateEndpointConnectionsListByResource(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "mhsmlistprivateendpointconnections.MHSMListPrivateEndpointConnectionsClient", "MHSMPrivateEndpointConnectionsListByResource", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
