package volumegroups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByElasticSanOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VolumeGroup
}

type ListByElasticSanCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VolumeGroup
}

// ListByElasticSan ...
func (c VolumeGroupsClient) ListByElasticSan(ctx context.Context, id ElasticSanId) (result ListByElasticSanOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/volumeGroups", id.ID()),
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
		Values *[]VolumeGroup `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByElasticSanComplete retrieves all the results into a single object
func (c VolumeGroupsClient) ListByElasticSanComplete(ctx context.Context, id ElasticSanId) (ListByElasticSanCompleteResult, error) {
	return c.ListByElasticSanCompleteMatchingPredicate(ctx, id, VolumeGroupOperationPredicate{})
}

// ListByElasticSanCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VolumeGroupsClient) ListByElasticSanCompleteMatchingPredicate(ctx context.Context, id ElasticSanId, predicate VolumeGroupOperationPredicate) (result ListByElasticSanCompleteResult, err error) {
	items := make([]VolumeGroup, 0)

	resp, err := c.ListByElasticSan(ctx, id)
	if err != nil {
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

	result = ListByElasticSanCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
