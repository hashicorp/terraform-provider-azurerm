package redis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedServerListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RedisLinkedServerWithProperties
}

type LinkedServerListCompleteResult struct {
	Items []RedisLinkedServerWithProperties
}

// LinkedServerList ...
func (c RedisClient) LinkedServerList(ctx context.Context, id RediId) (result LinkedServerListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/linkedServers", id.ID()),
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
		Values *[]RedisLinkedServerWithProperties `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LinkedServerListComplete retrieves all the results into a single object
func (c RedisClient) LinkedServerListComplete(ctx context.Context, id RediId) (LinkedServerListCompleteResult, error) {
	return c.LinkedServerListCompleteMatchingPredicate(ctx, id, RedisLinkedServerWithPropertiesOperationPredicate{})
}

// LinkedServerListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RedisClient) LinkedServerListCompleteMatchingPredicate(ctx context.Context, id RediId, predicate RedisLinkedServerWithPropertiesOperationPredicate) (result LinkedServerListCompleteResult, err error) {
	items := make([]RedisLinkedServerWithProperties, 0)

	resp, err := c.LinkedServerList(ctx, id)
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

	result = LinkedServerListCompleteResult{
		Items: items,
	}
	return
}
