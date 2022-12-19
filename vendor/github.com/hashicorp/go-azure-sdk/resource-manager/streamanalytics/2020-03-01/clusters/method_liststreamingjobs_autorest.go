package clusters

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

type ListStreamingJobsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ClusterJob

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListStreamingJobsOperationResponse, error)
}

type ListStreamingJobsCompleteResult struct {
	Items []ClusterJob
}

func (r ListStreamingJobsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListStreamingJobsOperationResponse) LoadMore(ctx context.Context) (resp ListStreamingJobsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListStreamingJobs ...
func (c ClustersClient) ListStreamingJobs(ctx context.Context, id ClusterId) (resp ListStreamingJobsOperationResponse, err error) {
	req, err := c.preparerForListStreamingJobs(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "ListStreamingJobs", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "ListStreamingJobs", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListStreamingJobs(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "ListStreamingJobs", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListStreamingJobs prepares the ListStreamingJobs request.
func (c ClustersClient) preparerForListStreamingJobs(ctx context.Context, id ClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listStreamingJobs", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListStreamingJobsWithNextLink prepares the ListStreamingJobs request with the given nextLink token.
func (c ClustersClient) preparerForListStreamingJobsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListStreamingJobs handles the response to the ListStreamingJobs request. The method always
// closes the http.Response Body.
func (c ClustersClient) responderForListStreamingJobs(resp *http.Response) (result ListStreamingJobsOperationResponse, err error) {
	type page struct {
		Values   []ClusterJob `json:"value"`
		NextLink *string      `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListStreamingJobsOperationResponse, err error) {
			req, err := c.preparerForListStreamingJobsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "ListStreamingJobs", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "ListStreamingJobs", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListStreamingJobs(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "ListStreamingJobs", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListStreamingJobsComplete retrieves all of the results into a single object
func (c ClustersClient) ListStreamingJobsComplete(ctx context.Context, id ClusterId) (ListStreamingJobsCompleteResult, error) {
	return c.ListStreamingJobsCompleteMatchingPredicate(ctx, id, ClusterJobOperationPredicate{})
}

// ListStreamingJobsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ClustersClient) ListStreamingJobsCompleteMatchingPredicate(ctx context.Context, id ClusterId, predicate ClusterJobOperationPredicate) (resp ListStreamingJobsCompleteResult, err error) {
	items := make([]ClusterJob, 0)

	page, err := c.ListStreamingJobs(ctx, id)
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

	out := ListStreamingJobsCompleteResult{
		Items: items,
	}
	return out, nil
}
