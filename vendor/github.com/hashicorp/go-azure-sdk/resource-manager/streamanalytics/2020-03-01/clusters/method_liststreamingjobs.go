package clusters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListStreamingJobsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ClusterJob
}

type ListStreamingJobsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ClusterJob
}

type ListStreamingJobsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListStreamingJobsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListStreamingJobs ...
func (c ClustersClient) ListStreamingJobs(ctx context.Context, id ClusterId) (result ListStreamingJobsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &ListStreamingJobsCustomPager{},
		Path:       fmt.Sprintf("%s/listStreamingJobs", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values *[]ClusterJob `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListStreamingJobsComplete retrieves all the results into a single object
func (c ClustersClient) ListStreamingJobsComplete(ctx context.Context, id ClusterId) (ListStreamingJobsCompleteResult, error) {
	return c.ListStreamingJobsCompleteMatchingPredicate(ctx, id, ClusterJobOperationPredicate{})
}

// ListStreamingJobsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ClustersClient) ListStreamingJobsCompleteMatchingPredicate(ctx context.Context, id ClusterId, predicate ClusterJobOperationPredicate) (result ListStreamingJobsCompleteResult, err error) {
	items := make([]ClusterJob, 0)

	resp, err := c.ListStreamingJobs(ctx, id)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	result = ListStreamingJobsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
