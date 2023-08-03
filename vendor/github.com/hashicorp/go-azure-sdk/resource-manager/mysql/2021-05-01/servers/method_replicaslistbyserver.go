package servers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicasListByServerOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Server
}

type ReplicasListByServerCompleteResult struct {
	Items []Server
}

// ReplicasListByServer ...
func (c ServersClient) ReplicasListByServer(ctx context.Context, id FlexibleServerId) (result ReplicasListByServerOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/replicas", id.ID()),
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
		Values *[]Server `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ReplicasListByServerComplete retrieves all the results into a single object
func (c ServersClient) ReplicasListByServerComplete(ctx context.Context, id FlexibleServerId) (ReplicasListByServerCompleteResult, error) {
	return c.ReplicasListByServerCompleteMatchingPredicate(ctx, id, ServerOperationPredicate{})
}

// ReplicasListByServerCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ServersClient) ReplicasListByServerCompleteMatchingPredicate(ctx context.Context, id FlexibleServerId, predicate ServerOperationPredicate) (result ReplicasListByServerCompleteResult, err error) {
	items := make([]Server, 0)

	resp, err := c.ReplicasListByServer(ctx, id)
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

	result = ReplicasListByServerCompleteResult{
		Items: items,
	}
	return
}
