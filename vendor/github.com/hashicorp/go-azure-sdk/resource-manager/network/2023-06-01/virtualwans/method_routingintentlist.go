package virtualwans

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoutingIntentListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RoutingIntent
}

type RoutingIntentListCompleteResult struct {
	Items []RoutingIntent
}

// RoutingIntentList ...
func (c VirtualWANsClient) RoutingIntentList(ctx context.Context, id VirtualHubId) (result RoutingIntentListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/routingIntent", id.ID()),
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
		Values *[]RoutingIntent `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RoutingIntentListComplete retrieves all the results into a single object
func (c VirtualWANsClient) RoutingIntentListComplete(ctx context.Context, id VirtualHubId) (RoutingIntentListCompleteResult, error) {
	return c.RoutingIntentListCompleteMatchingPredicate(ctx, id, RoutingIntentOperationPredicate{})
}

// RoutingIntentListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) RoutingIntentListCompleteMatchingPredicate(ctx context.Context, id VirtualHubId, predicate RoutingIntentOperationPredicate) (result RoutingIntentListCompleteResult, err error) {
	items := make([]RoutingIntent, 0)

	resp, err := c.RoutingIntentList(ctx, id)
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

	result = RoutingIntentListCompleteResult{
		Items: items,
	}
	return
}
