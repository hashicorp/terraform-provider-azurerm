package firewallstatus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByFirewallsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FirewallStatusResource
}

type ListByFirewallsCompleteResult struct {
	Items []FirewallStatusResource
}

// ListByFirewalls ...
func (c FirewallStatusClient) ListByFirewalls(ctx context.Context, id FirewallId) (result ListByFirewallsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/statuses", id.ID()),
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
		Values *[]FirewallStatusResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByFirewallsComplete retrieves all the results into a single object
func (c FirewallStatusClient) ListByFirewallsComplete(ctx context.Context, id FirewallId) (ListByFirewallsCompleteResult, error) {
	return c.ListByFirewallsCompleteMatchingPredicate(ctx, id, FirewallStatusResourceOperationPredicate{})
}

// ListByFirewallsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c FirewallStatusClient) ListByFirewallsCompleteMatchingPredicate(ctx context.Context, id FirewallId, predicate FirewallStatusResourceOperationPredicate) (result ListByFirewallsCompleteResult, err error) {
	items := make([]FirewallStatusResource, 0)

	resp, err := c.ListByFirewalls(ctx, id)
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

	result = ListByFirewallsCompleteResult{
		Items: items,
	}
	return
}
