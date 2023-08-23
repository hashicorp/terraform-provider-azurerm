package healthbots

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

type BotsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]HealthBot
}

type BotsListByResourceGroupCompleteResult struct {
	Items []HealthBot
}

// BotsListByResourceGroup ...
func (c HealthbotsClient) BotsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result BotsListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.HealthBot/healthBots", id.ID()),
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
		Values *[]HealthBot `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// BotsListByResourceGroupComplete retrieves all the results into a single object
func (c HealthbotsClient) BotsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (BotsListByResourceGroupCompleteResult, error) {
	return c.BotsListByResourceGroupCompleteMatchingPredicate(ctx, id, HealthBotOperationPredicate{})
}

// BotsListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c HealthbotsClient) BotsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate HealthBotOperationPredicate) (result BotsListByResourceGroupCompleteResult, err error) {
	items := make([]HealthBot, 0)

	resp, err := c.BotsListByResourceGroup(ctx, id)
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

	result = BotsListByResourceGroupCompleteResult{
		Items: items,
	}
	return
}
