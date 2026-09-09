package appserviceplans

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

type ListHybridConnectionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]HybridConnection
}

type ListHybridConnectionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []HybridConnection
}

type ListHybridConnectionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListHybridConnectionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListHybridConnections ...
func (c AppServicePlansClient) ListHybridConnections(ctx context.Context, id commonids.AppServicePlanId) (result ListHybridConnectionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListHybridConnectionsCustomPager{},
		Path:       fmt.Sprintf("%s/hybridConnectionRelays", id.ID()),
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
		Values *[]HybridConnection `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListHybridConnectionsComplete retrieves all the results into a single object
func (c AppServicePlansClient) ListHybridConnectionsComplete(ctx context.Context, id commonids.AppServicePlanId) (ListHybridConnectionsCompleteResult, error) {
	return c.ListHybridConnectionsCompleteMatchingPredicate(ctx, id, HybridConnectionOperationPredicate{})
}

// ListHybridConnectionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppServicePlansClient) ListHybridConnectionsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServicePlanId, predicate HybridConnectionOperationPredicate) (result ListHybridConnectionsCompleteResult, err error) {
	items := make([]HybridConnection, 0)

	resp, err := c.ListHybridConnections(ctx, id)
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

	result = ListHybridConnectionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
