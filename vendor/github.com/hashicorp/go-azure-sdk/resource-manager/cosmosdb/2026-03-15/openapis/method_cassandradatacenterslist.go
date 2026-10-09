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

type CassandraDataCentersListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DataCenterResource
}

type CassandraDataCentersListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DataCenterResource
}

type CassandraDataCentersListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CassandraDataCentersListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CassandraDataCentersList ...
func (c OpenapisClient) CassandraDataCentersList(ctx context.Context, id CassandraClusterId) (result CassandraDataCentersListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CassandraDataCentersListCustomPager{},
		Path:       fmt.Sprintf("%s/dataCenters", id.ID()),
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
		Values *[]DataCenterResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CassandraDataCentersListComplete retrieves all the results into a single object
func (c OpenapisClient) CassandraDataCentersListComplete(ctx context.Context, id CassandraClusterId) (CassandraDataCentersListCompleteResult, error) {
	return c.CassandraDataCentersListCompleteMatchingPredicate(ctx, id, DataCenterResourceOperationPredicate{})
}

// CassandraDataCentersListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) CassandraDataCentersListCompleteMatchingPredicate(ctx context.Context, id CassandraClusterId, predicate DataCenterResourceOperationPredicate) (result CassandraDataCentersListCompleteResult, err error) {
	items := make([]DataCenterResource, 0)

	resp, err := c.CassandraDataCentersList(ctx, id)
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

	result = CassandraDataCentersListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
