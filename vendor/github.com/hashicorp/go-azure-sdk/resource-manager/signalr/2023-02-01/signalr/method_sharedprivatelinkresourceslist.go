package signalr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharedPrivateLinkResourcesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SharedPrivateLinkResource
}

type SharedPrivateLinkResourcesListCompleteResult struct {
	Items []SharedPrivateLinkResource
}

// SharedPrivateLinkResourcesList ...
func (c SignalRClient) SharedPrivateLinkResourcesList(ctx context.Context, id SignalRId) (result SharedPrivateLinkResourcesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/sharedPrivateLinkResources", id.ID()),
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
		Values *[]SharedPrivateLinkResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SharedPrivateLinkResourcesListComplete retrieves all the results into a single object
func (c SignalRClient) SharedPrivateLinkResourcesListComplete(ctx context.Context, id SignalRId) (SharedPrivateLinkResourcesListCompleteResult, error) {
	return c.SharedPrivateLinkResourcesListCompleteMatchingPredicate(ctx, id, SharedPrivateLinkResourceOperationPredicate{})
}

// SharedPrivateLinkResourcesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SignalRClient) SharedPrivateLinkResourcesListCompleteMatchingPredicate(ctx context.Context, id SignalRId, predicate SharedPrivateLinkResourceOperationPredicate) (result SharedPrivateLinkResourcesListCompleteResult, err error) {
	items := make([]SharedPrivateLinkResource, 0)

	resp, err := c.SharedPrivateLinkResourcesList(ctx, id)
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

	result = SharedPrivateLinkResourcesListCompleteResult{
		Items: items,
	}
	return
}
