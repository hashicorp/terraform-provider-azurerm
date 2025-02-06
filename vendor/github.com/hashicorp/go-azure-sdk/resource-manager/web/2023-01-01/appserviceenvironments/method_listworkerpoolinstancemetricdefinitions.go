package appserviceenvironments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListWorkerPoolInstanceMetricDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ResourceMetricDefinition
}

type ListWorkerPoolInstanceMetricDefinitionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ResourceMetricDefinition
}

type ListWorkerPoolInstanceMetricDefinitionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListWorkerPoolInstanceMetricDefinitionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListWorkerPoolInstanceMetricDefinitions ...
func (c AppServiceEnvironmentsClient) ListWorkerPoolInstanceMetricDefinitions(ctx context.Context, id WorkerPoolInstanceId) (result ListWorkerPoolInstanceMetricDefinitionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListWorkerPoolInstanceMetricDefinitionsCustomPager{},
		Path:       fmt.Sprintf("%s/metricdefinitions", id.ID()),
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
		Values *[]ResourceMetricDefinition `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListWorkerPoolInstanceMetricDefinitionsComplete retrieves all the results into a single object
func (c AppServiceEnvironmentsClient) ListWorkerPoolInstanceMetricDefinitionsComplete(ctx context.Context, id WorkerPoolInstanceId) (ListWorkerPoolInstanceMetricDefinitionsCompleteResult, error) {
	return c.ListWorkerPoolInstanceMetricDefinitionsCompleteMatchingPredicate(ctx, id, ResourceMetricDefinitionOperationPredicate{})
}

// ListWorkerPoolInstanceMetricDefinitionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppServiceEnvironmentsClient) ListWorkerPoolInstanceMetricDefinitionsCompleteMatchingPredicate(ctx context.Context, id WorkerPoolInstanceId, predicate ResourceMetricDefinitionOperationPredicate) (result ListWorkerPoolInstanceMetricDefinitionsCompleteResult, err error) {
	items := make([]ResourceMetricDefinition, 0)

	resp, err := c.ListWorkerPoolInstanceMetricDefinitions(ctx, id)
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

	result = ListWorkerPoolInstanceMetricDefinitionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
