package querypacks

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

type QueryPacksListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LogAnalyticsQueryPack
}

type QueryPacksListByResourceGroupCompleteResult struct {
	Items []LogAnalyticsQueryPack
}

// QueryPacksListByResourceGroup ...
func (c QueryPacksClient) QueryPacksListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result QueryPacksListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.OperationalInsights/queryPacks", id.ID()),
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
		Values *[]LogAnalyticsQueryPack `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// QueryPacksListByResourceGroupComplete retrieves all the results into a single object
func (c QueryPacksClient) QueryPacksListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (QueryPacksListByResourceGroupCompleteResult, error) {
	return c.QueryPacksListByResourceGroupCompleteMatchingPredicate(ctx, id, LogAnalyticsQueryPackOperationPredicate{})
}

// QueryPacksListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c QueryPacksClient) QueryPacksListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate LogAnalyticsQueryPackOperationPredicate) (result QueryPacksListByResourceGroupCompleteResult, err error) {
	items := make([]LogAnalyticsQueryPack, 0)

	resp, err := c.QueryPacksListByResourceGroup(ctx, id)
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

	result = QueryPacksListByResourceGroupCompleteResult{
		Items: items,
	}
	return
}
