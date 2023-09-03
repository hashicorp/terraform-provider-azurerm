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

type QueryPacksListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LogAnalyticsQueryPack
}

type QueryPacksListCompleteResult struct {
	Items []LogAnalyticsQueryPack
}

// QueryPacksList ...
func (c QueryPacksClient) QueryPacksList(ctx context.Context, id commonids.SubscriptionId) (result QueryPacksListOperationResponse, err error) {
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

// QueryPacksListComplete retrieves all the results into a single object
func (c QueryPacksClient) QueryPacksListComplete(ctx context.Context, id commonids.SubscriptionId) (QueryPacksListCompleteResult, error) {
	return c.QueryPacksListCompleteMatchingPredicate(ctx, id, LogAnalyticsQueryPackOperationPredicate{})
}

// QueryPacksListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c QueryPacksClient) QueryPacksListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate LogAnalyticsQueryPackOperationPredicate) (result QueryPacksListCompleteResult, err error) {
	items := make([]LogAnalyticsQueryPack, 0)

	resp, err := c.QueryPacksList(ctx, id)
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

	result = QueryPacksListCompleteResult{
		Items: items,
	}
	return
}
