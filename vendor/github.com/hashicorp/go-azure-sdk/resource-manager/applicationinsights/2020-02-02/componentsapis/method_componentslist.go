package componentsapis

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

type ComponentsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApplicationInsightsComponent
}

type ComponentsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ApplicationInsightsComponent
}

// ComponentsList ...
func (c ComponentsAPIsClient) ComponentsList(ctx context.Context, id commonids.SubscriptionId) (result ComponentsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.Insights/components", id.ID()),
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
		Values *[]ApplicationInsightsComponent `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ComponentsListComplete retrieves all the results into a single object
func (c ComponentsAPIsClient) ComponentsListComplete(ctx context.Context, id commonids.SubscriptionId) (ComponentsListCompleteResult, error) {
	return c.ComponentsListCompleteMatchingPredicate(ctx, id, ApplicationInsightsComponentOperationPredicate{})
}

// ComponentsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ComponentsAPIsClient) ComponentsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ApplicationInsightsComponentOperationPredicate) (result ComponentsListCompleteResult, err error) {
	items := make([]ApplicationInsightsComponent, 0)

	resp, err := c.ComponentsList(ctx, id)
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

	result = ComponentsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
