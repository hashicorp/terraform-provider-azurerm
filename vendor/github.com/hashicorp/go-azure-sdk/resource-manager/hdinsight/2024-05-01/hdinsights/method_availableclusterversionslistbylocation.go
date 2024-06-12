package hdinsights

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableClusterVersionsListByLocationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ClusterVersion
}

type AvailableClusterVersionsListByLocationCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ClusterVersion
}

// AvailableClusterVersionsListByLocation ...
func (c HdinsightsClient) AvailableClusterVersionsListByLocation(ctx context.Context, id LocationId) (result AvailableClusterVersionsListByLocationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/availableClusterVersions", id.ID()),
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
		Values *[]ClusterVersion `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AvailableClusterVersionsListByLocationComplete retrieves all the results into a single object
func (c HdinsightsClient) AvailableClusterVersionsListByLocationComplete(ctx context.Context, id LocationId) (AvailableClusterVersionsListByLocationCompleteResult, error) {
	return c.AvailableClusterVersionsListByLocationCompleteMatchingPredicate(ctx, id, ClusterVersionOperationPredicate{})
}

// AvailableClusterVersionsListByLocationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c HdinsightsClient) AvailableClusterVersionsListByLocationCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate ClusterVersionOperationPredicate) (result AvailableClusterVersionsListByLocationCompleteResult, err error) {
	items := make([]ClusterVersion, 0)

	resp, err := c.AvailableClusterVersionsListByLocation(ctx, id)
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

	result = AvailableClusterVersionsListByLocationCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
