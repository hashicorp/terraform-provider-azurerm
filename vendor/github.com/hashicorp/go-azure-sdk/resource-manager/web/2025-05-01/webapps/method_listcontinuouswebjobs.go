package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListContinuousWebJobsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ContinuousWebJob
}

type ListContinuousWebJobsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ContinuousWebJob
}

type ListContinuousWebJobsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListContinuousWebJobsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListContinuousWebJobs ...
func (c WebAppsClient) ListContinuousWebJobs(ctx context.Context, id commonids.AppServiceId) (result ListContinuousWebJobsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListContinuousWebJobsCustomPager{},
		Path:       fmt.Sprintf("%s/continuousWebJobs", id.ID()),
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
		Values *[]ContinuousWebJob `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListContinuousWebJobsComplete retrieves all the results into a single object
func (c WebAppsClient) ListContinuousWebJobsComplete(ctx context.Context, id commonids.AppServiceId) (ListContinuousWebJobsCompleteResult, error) {
	return c.ListContinuousWebJobsCompleteMatchingPredicate(ctx, id, ContinuousWebJobOperationPredicate{})
}

// ListContinuousWebJobsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListContinuousWebJobsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate ContinuousWebJobOperationPredicate) (result ListContinuousWebJobsCompleteResult, err error) {
	items := make([]ContinuousWebJob, 0)

	resp, err := c.ListContinuousWebJobs(ctx, id)
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

	result = ListContinuousWebJobsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
