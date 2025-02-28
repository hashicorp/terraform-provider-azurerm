package apimanagementgatewayskus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAvailableSkusOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GatewayResourceSkuResult
}

type ListAvailableSkusCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GatewayResourceSkuResult
}

type ListAvailableSkusCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAvailableSkusCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAvailableSkus ...
func (c ApiManagementGatewaySkusClient) ListAvailableSkus(ctx context.Context, id GatewayId) (result ListAvailableSkusOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListAvailableSkusCustomPager{},
		Path:       fmt.Sprintf("%s/skus", id.ID()),
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
		Values *[]GatewayResourceSkuResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAvailableSkusComplete retrieves all the results into a single object
func (c ApiManagementGatewaySkusClient) ListAvailableSkusComplete(ctx context.Context, id GatewayId) (ListAvailableSkusCompleteResult, error) {
	return c.ListAvailableSkusCompleteMatchingPredicate(ctx, id, GatewayResourceSkuResultOperationPredicate{})
}

// ListAvailableSkusCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApiManagementGatewaySkusClient) ListAvailableSkusCompleteMatchingPredicate(ctx context.Context, id GatewayId, predicate GatewayResourceSkuResultOperationPredicate) (result ListAvailableSkusCompleteResult, err error) {
	items := make([]GatewayResourceSkuResult, 0)

	resp, err := c.ListAvailableSkus(ctx, id)
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

	result = ListAvailableSkusCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
