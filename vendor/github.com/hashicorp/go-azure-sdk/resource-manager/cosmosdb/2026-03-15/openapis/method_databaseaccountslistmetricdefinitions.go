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

type DatabaseAccountsListMetricDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MetricDefinition
}

type DatabaseAccountsListMetricDefinitionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []MetricDefinition
}

type DatabaseAccountsListMetricDefinitionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DatabaseAccountsListMetricDefinitionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DatabaseAccountsListMetricDefinitions ...
func (c OpenapisClient) DatabaseAccountsListMetricDefinitions(ctx context.Context, id DatabaseAccountId) (result DatabaseAccountsListMetricDefinitionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DatabaseAccountsListMetricDefinitionsCustomPager{},
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

// DatabaseAccountsListMetricDefinitionsComplete retrieves all the results into a single object
func (c OpenapisClient) DatabaseAccountsListMetricDefinitionsComplete(ctx context.Context, id DatabaseAccountId) (DatabaseAccountsListMetricDefinitionsCompleteResult, error) {
	return c.DatabaseAccountsListMetricDefinitionsCompleteMatchingPredicate(ctx, id, MetricDefinitionOperationPredicate{})
}

// DatabaseAccountsListMetricDefinitionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) DatabaseAccountsListMetricDefinitionsCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate MetricDefinitionOperationPredicate) (result DatabaseAccountsListMetricDefinitionsCompleteResult, err error) {
	items := make([]MetricDefinition, 0)

	resp, err := c.DatabaseAccountsListMetricDefinitions(ctx, id)
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

	result = DatabaseAccountsListMetricDefinitionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
