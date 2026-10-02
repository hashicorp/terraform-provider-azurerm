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

type CollectionListMetricDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MetricDefinition
}

type CollectionListMetricDefinitionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []MetricDefinition
}

type CollectionListMetricDefinitionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CollectionListMetricDefinitionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CollectionListMetricDefinitions ...
func (c OpenapisClient) CollectionListMetricDefinitions(ctx context.Context, id CollectionId) (result CollectionListMetricDefinitionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CollectionListMetricDefinitionsCustomPager{},
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

// CollectionListMetricDefinitionsComplete retrieves all the results into a single object
func (c OpenapisClient) CollectionListMetricDefinitionsComplete(ctx context.Context, id CollectionId) (CollectionListMetricDefinitionsCompleteResult, error) {
	return c.CollectionListMetricDefinitionsCompleteMatchingPredicate(ctx, id, MetricDefinitionOperationPredicate{})
}

// CollectionListMetricDefinitionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) CollectionListMetricDefinitionsCompleteMatchingPredicate(ctx context.Context, id CollectionId, predicate MetricDefinitionOperationPredicate) (result CollectionListMetricDefinitionsCompleteResult, err error) {
	items := make([]MetricDefinition, 0)

	resp, err := c.CollectionListMetricDefinitions(ctx, id)
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

	result = CollectionListMetricDefinitionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
