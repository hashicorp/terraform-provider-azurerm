package openapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseListMetricDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MetricDefinition
}

type DatabaseListMetricDefinitionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []MetricDefinition
}

type DatabaseListMetricDefinitionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DatabaseListMetricDefinitionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DatabaseListMetricDefinitions ...
func (c OpenapisClient) DatabaseListMetricDefinitions(ctx context.Context, id DatabaseId) (result DatabaseListMetricDefinitionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DatabaseListMetricDefinitionsCustomPager{},
		Path:       fmt.Sprintf("%s/metricDefinitions", id.ID()),
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
		Values *[]MetricDefinition `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DatabaseListMetricDefinitionsComplete retrieves all the results into a single object
func (c OpenapisClient) DatabaseListMetricDefinitionsComplete(ctx context.Context, id DatabaseId) (DatabaseListMetricDefinitionsCompleteResult, error) {
	return c.DatabaseListMetricDefinitionsCompleteMatchingPredicate(ctx, id, MetricDefinitionOperationPredicate{})
}

// DatabaseListMetricDefinitionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) DatabaseListMetricDefinitionsCompleteMatchingPredicate(ctx context.Context, id DatabaseId, predicate MetricDefinitionOperationPredicate) (result DatabaseListMetricDefinitionsCompleteResult, err error) {
	items := make([]MetricDefinition, 0)

	resp, err := c.DatabaseListMetricDefinitions(ctx, id)
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

	result = DatabaseListMetricDefinitionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
