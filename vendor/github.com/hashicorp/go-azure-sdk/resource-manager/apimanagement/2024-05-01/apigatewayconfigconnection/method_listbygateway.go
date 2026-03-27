package apigatewayconfigconnection

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByGatewayOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApiManagementGatewayConfigConnectionResource
}

type ListByGatewayCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ApiManagementGatewayConfigConnectionResource
}

type ListByGatewayCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByGatewayCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByGateway ...
func (c ApiGatewayConfigConnectionClient) ListByGateway(ctx context.Context, id GatewayId) (result ListByGatewayOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByGatewayCustomPager{},
		Path:       fmt.Sprintf("%s/configConnections", id.ID()),
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
		Values *[]ApiManagementGatewayConfigConnectionResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByGatewayComplete retrieves all the results into a single object
func (c ApiGatewayConfigConnectionClient) ListByGatewayComplete(ctx context.Context, id GatewayId) (ListByGatewayCompleteResult, error) {
	return c.ListByGatewayCompleteMatchingPredicate(ctx, id, ApiManagementGatewayConfigConnectionResourceOperationPredicate{})
}

// ListByGatewayCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApiGatewayConfigConnectionClient) ListByGatewayCompleteMatchingPredicate(ctx context.Context, id GatewayId, predicate ApiManagementGatewayConfigConnectionResourceOperationPredicate) (result ListByGatewayCompleteResult, err error) {
	items := make([]ApiManagementGatewayConfigConnectionResource, 0)

	resp, err := c.ListByGateway(ctx, id)
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

	result = ListByGatewayCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
