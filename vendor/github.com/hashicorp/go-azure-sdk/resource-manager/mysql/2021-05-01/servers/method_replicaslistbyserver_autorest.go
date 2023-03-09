package servers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicasListByServerOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Server

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ReplicasListByServerOperationResponse, error)
}

type ReplicasListByServerCompleteResult struct {
	Items []Server
}

func (r ReplicasListByServerOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ReplicasListByServerOperationResponse) LoadMore(ctx context.Context) (resp ReplicasListByServerOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ReplicasListByServer ...
func (c ServersClient) ReplicasListByServer(ctx context.Context, id FlexibleServerId) (resp ReplicasListByServerOperationResponse, err error) {
	req, err := c.preparerForReplicasListByServer(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servers.ServersClient", "ReplicasListByServer", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "servers.ServersClient", "ReplicasListByServer", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForReplicasListByServer(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servers.ServersClient", "ReplicasListByServer", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForReplicasListByServer prepares the ReplicasListByServer request.
func (c ServersClient) preparerForReplicasListByServer(ctx context.Context, id FlexibleServerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/replicas", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForReplicasListByServerWithNextLink prepares the ReplicasListByServer request with the given nextLink token.
func (c ServersClient) preparerForReplicasListByServerWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForReplicasListByServer handles the response to the ReplicasListByServer request. The method always
// closes the http.Response Body.
func (c ServersClient) responderForReplicasListByServer(resp *http.Response) (result ReplicasListByServerOperationResponse, err error) {
	type page struct {
		Values   []Server `json:"value"`
		NextLink *string  `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ReplicasListByServerOperationResponse, err error) {
			req, err := c.preparerForReplicasListByServerWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "servers.ServersClient", "ReplicasListByServer", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "servers.ServersClient", "ReplicasListByServer", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForReplicasListByServer(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "servers.ServersClient", "ReplicasListByServer", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ReplicasListByServerComplete retrieves all of the results into a single object
func (c ServersClient) ReplicasListByServerComplete(ctx context.Context, id FlexibleServerId) (ReplicasListByServerCompleteResult, error) {
	return c.ReplicasListByServerCompleteMatchingPredicate(ctx, id, ServerOperationPredicate{})
}

// ReplicasListByServerCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ServersClient) ReplicasListByServerCompleteMatchingPredicate(ctx context.Context, id FlexibleServerId, predicate ServerOperationPredicate) (resp ReplicasListByServerCompleteResult, err error) {
	items := make([]Server, 0)

	page, err := c.ReplicasListByServer(ctx, id)
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

	out := ReplicasListByServerCompleteResult{
		Items: items,
	}
	return out, nil
}
