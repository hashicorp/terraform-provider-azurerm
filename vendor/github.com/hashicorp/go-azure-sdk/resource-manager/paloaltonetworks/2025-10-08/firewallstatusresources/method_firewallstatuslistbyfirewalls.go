package firewallstatusresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallStatusListByFirewallsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FirewallStatusResource
}

type FirewallStatusListByFirewallsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []FirewallStatusResource
}

type FirewallStatusListByFirewallsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *FirewallStatusListByFirewallsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// FirewallStatusListByFirewalls ...
func (c FirewallStatusResourcesClient) FirewallStatusListByFirewalls(ctx context.Context, id FirewallId) (result FirewallStatusListByFirewallsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &FirewallStatusListByFirewallsCustomPager{},
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

// FirewallStatusListByFirewallsComplete retrieves all the results into a single object
func (c FirewallStatusResourcesClient) FirewallStatusListByFirewallsComplete(ctx context.Context, id FirewallId) (FirewallStatusListByFirewallsCompleteResult, error) {
	return c.FirewallStatusListByFirewallsCompleteMatchingPredicate(ctx, id, FirewallStatusResourceOperationPredicate{})
}

// FirewallStatusListByFirewallsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c FirewallStatusResourcesClient) FirewallStatusListByFirewallsCompleteMatchingPredicate(ctx context.Context, id FirewallId, predicate FirewallStatusResourceOperationPredicate) (result FirewallStatusListByFirewallsCompleteResult, err error) {
	items := make([]FirewallStatusResource, 0)

	resp, err := c.FirewallStatusListByFirewalls(ctx, id)
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

	result = FirewallStatusListByFirewallsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
