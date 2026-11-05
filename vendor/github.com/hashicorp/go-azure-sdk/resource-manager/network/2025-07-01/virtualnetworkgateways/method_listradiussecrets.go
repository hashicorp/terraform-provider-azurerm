package virtualnetworkgateways

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListRadiusSecretsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RadiusAuthServer
}

type ListRadiusSecretsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RadiusAuthServer
}

type ListRadiusSecretsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListRadiusSecretsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListRadiusSecrets ...
func (c VirtualNetworkGatewaysClient) ListRadiusSecrets(ctx context.Context, id VirtualNetworkGatewayId) (result ListRadiusSecretsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &ListRadiusSecretsCustomPager{},
		Path:       fmt.Sprintf("%s/listRadiusSecrets", id.ID()),
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
		Values *[]RadiusAuthServer `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListRadiusSecretsComplete retrieves all the results into a single object
func (c VirtualNetworkGatewaysClient) ListRadiusSecretsComplete(ctx context.Context, id VirtualNetworkGatewayId) (ListRadiusSecretsCompleteResult, error) {
	return c.ListRadiusSecretsCompleteMatchingPredicate(ctx, id, RadiusAuthServerOperationPredicate{})
}

// ListRadiusSecretsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualNetworkGatewaysClient) ListRadiusSecretsCompleteMatchingPredicate(ctx context.Context, id VirtualNetworkGatewayId, predicate RadiusAuthServerOperationPredicate) (result ListRadiusSecretsCompleteResult, err error) {
	items := make([]RadiusAuthServer, 0)

	resp, err := c.ListRadiusSecrets(ctx, id)
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

	result = ListRadiusSecretsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
