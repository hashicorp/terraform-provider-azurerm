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

type ListMultiRolePoolInstanceMetricDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ResourceMetricDefinition
}

type ListMultiRolePoolInstanceMetricDefinitionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ResourceMetricDefinition
}

type ListMultiRolePoolInstanceMetricDefinitionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListMultiRolePoolInstanceMetricDefinitionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListMultiRolePoolInstanceMetricDefinitions ...
func (c AppServiceEnvironmentsClient) ListMultiRolePoolInstanceMetricDefinitions(ctx context.Context, id DefaultInstanceId) (result ListMultiRolePoolInstanceMetricDefinitionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListMultiRolePoolInstanceMetricDefinitionsCustomPager{},
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

// ListMultiRolePoolInstanceMetricDefinitionsComplete retrieves all the results into a single object
func (c AppServiceEnvironmentsClient) ListMultiRolePoolInstanceMetricDefinitionsComplete(ctx context.Context, id DefaultInstanceId) (ListMultiRolePoolInstanceMetricDefinitionsCompleteResult, error) {
	return c.ListMultiRolePoolInstanceMetricDefinitionsCompleteMatchingPredicate(ctx, id, ResourceMetricDefinitionOperationPredicate{})
}

// ListMultiRolePoolInstanceMetricDefinitionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppServiceEnvironmentsClient) ListMultiRolePoolInstanceMetricDefinitionsCompleteMatchingPredicate(ctx context.Context, id DefaultInstanceId, predicate ResourceMetricDefinitionOperationPredicate) (result ListMultiRolePoolInstanceMetricDefinitionsCompleteResult, err error) {
	items := make([]ResourceMetricDefinition, 0)

	resp, err := c.ListMultiRolePoolInstanceMetricDefinitions(ctx, id)
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

	result = ListMultiRolePoolInstanceMetricDefinitionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
