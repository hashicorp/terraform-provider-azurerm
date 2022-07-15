package machinelearningcomputes

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

type ComputeListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ComputeResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ComputeListOperationResponse, error)
}

type ComputeListCompleteResult struct {
	Items []ComputeResource
}

func (r ComputeListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ComputeListOperationResponse) LoadMore(ctx context.Context) (resp ComputeListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ComputeListOperationOptions struct {
	Skip *string
}

func DefaultComputeListOperationOptions() ComputeListOperationOptions {
	return ComputeListOperationOptions{}
}

func (o ComputeListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ComputeListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Skip != nil {
		out["$skip"] = *o.Skip
	}

	return out
}

// ComputeList ...
func (c MachineLearningComputesClient) ComputeList(ctx context.Context, id WorkspaceId, options ComputeListOperationOptions) (resp ComputeListOperationResponse, err error) {
	req, err := c.preparerForComputeList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForComputeList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ComputeListComplete retrieves all of the results into a single object
func (c MachineLearningComputesClient) ComputeListComplete(ctx context.Context, id WorkspaceId, options ComputeListOperationOptions) (ComputeListCompleteResult, error) {
	return c.ComputeListCompleteMatchingPredicate(ctx, id, options, ComputeResourceOperationPredicate{})
}

// ComputeListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c MachineLearningComputesClient) ComputeListCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options ComputeListOperationOptions, predicate ComputeResourceOperationPredicate) (resp ComputeListCompleteResult, err error) {
	items := make([]ComputeResource, 0)

	page, err := c.ComputeList(ctx, id, options)
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

	out := ComputeListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForComputeList prepares the ComputeList request.
func (c MachineLearningComputesClient) preparerForComputeList(ctx context.Context, id WorkspaceId, options ComputeListOperationOptions) (*http.Request, error) {
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
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/computes", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForComputeListWithNextLink prepares the ComputeList request with the given nextLink token.
func (c MachineLearningComputesClient) preparerForComputeListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForComputeList handles the response to the ComputeList request. The method always
// closes the http.Response Body.
func (c MachineLearningComputesClient) responderForComputeList(resp *http.Response) (result ComputeListOperationResponse, err error) {
	type page struct {
		Values   []ComputeResource `json:"value"`
		NextLink *string           `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ComputeListOperationResponse, err error) {
			req, err := c.preparerForComputeListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForComputeList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "machinelearningcomputes.MachineLearningComputesClient", "ComputeList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
