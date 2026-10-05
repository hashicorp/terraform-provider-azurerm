package appserviceenvironments

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

type ListWorkerPoolsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]WorkerPoolResource
}

type ListWorkerPoolsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []WorkerPoolResource
}

type ListWorkerPoolsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListWorkerPoolsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListWorkerPools ...
func (c AppServiceEnvironmentsClient) ListWorkerPools(ctx context.Context, id commonids.AppServiceEnvironmentId) (result ListWorkerPoolsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListWorkerPoolsCustomPager{},
		Path:       fmt.Sprintf("%s/workerPools", id.ID()),
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
		Values *[]WorkerPoolResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListWorkerPoolsComplete retrieves all the results into a single object
func (c AppServiceEnvironmentsClient) ListWorkerPoolsComplete(ctx context.Context, id commonids.AppServiceEnvironmentId) (ListWorkerPoolsCompleteResult, error) {
	return c.ListWorkerPoolsCompleteMatchingPredicate(ctx, id, WorkerPoolResourceOperationPredicate{})
}

// ListWorkerPoolsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppServiceEnvironmentsClient) ListWorkerPoolsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceEnvironmentId, predicate WorkerPoolResourceOperationPredicate) (result ListWorkerPoolsCompleteResult, err error) {
	items := make([]WorkerPoolResource, 0)

	resp, err := c.ListWorkerPools(ctx, id)
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

	result = ListWorkerPoolsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
