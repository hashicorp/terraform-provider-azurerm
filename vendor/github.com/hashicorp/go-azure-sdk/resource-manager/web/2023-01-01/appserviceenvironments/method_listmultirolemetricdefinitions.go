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

type ListMultiRoleMetricDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ResourceMetricDefinition
}

type ListMultiRoleMetricDefinitionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ResourceMetricDefinition
}

type ListMultiRoleMetricDefinitionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListMultiRoleMetricDefinitionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListMultiRoleMetricDefinitions ...
func (c AppServiceEnvironmentsClient) ListMultiRoleMetricDefinitions(ctx context.Context, id commonids.AppServiceEnvironmentId) (result ListMultiRoleMetricDefinitionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListMultiRoleMetricDefinitionsCustomPager{},
		Path:       fmt.Sprintf("%s/multiRolePools/default/metricdefinitions", id.ID()),
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

// ListMultiRoleMetricDefinitionsComplete retrieves all the results into a single object
func (c AppServiceEnvironmentsClient) ListMultiRoleMetricDefinitionsComplete(ctx context.Context, id commonids.AppServiceEnvironmentId) (ListMultiRoleMetricDefinitionsCompleteResult, error) {
	return c.ListMultiRoleMetricDefinitionsCompleteMatchingPredicate(ctx, id, ResourceMetricDefinitionOperationPredicate{})
}

// ListMultiRoleMetricDefinitionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppServiceEnvironmentsClient) ListMultiRoleMetricDefinitionsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceEnvironmentId, predicate ResourceMetricDefinitionOperationPredicate) (result ListMultiRoleMetricDefinitionsCompleteResult, err error) {
	items := make([]ResourceMetricDefinition, 0)

	resp, err := c.ListMultiRoleMetricDefinitions(ctx, id)
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

	result = ListMultiRoleMetricDefinitionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
