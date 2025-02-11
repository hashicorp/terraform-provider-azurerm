package testlines

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByCommunicationsGatewayOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TestLine
}

type ListByCommunicationsGatewayCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TestLine
}

type ListByCommunicationsGatewayCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByCommunicationsGatewayCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByCommunicationsGateway ...
func (c TestLinesClient) ListByCommunicationsGateway(ctx context.Context, id CommunicationsGatewayId) (result ListByCommunicationsGatewayOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByCommunicationsGatewayCustomPager{},
		Path:       fmt.Sprintf("%s/testLines", id.ID()),
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
		Values *[]TestLine `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByCommunicationsGatewayComplete retrieves all the results into a single object
func (c TestLinesClient) ListByCommunicationsGatewayComplete(ctx context.Context, id CommunicationsGatewayId) (ListByCommunicationsGatewayCompleteResult, error) {
	return c.ListByCommunicationsGatewayCompleteMatchingPredicate(ctx, id, TestLineOperationPredicate{})
}

// ListByCommunicationsGatewayCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c TestLinesClient) ListByCommunicationsGatewayCompleteMatchingPredicate(ctx context.Context, id CommunicationsGatewayId, predicate TestLineOperationPredicate) (result ListByCommunicationsGatewayCompleteResult, err error) {
	items := make([]TestLine, 0)

	resp, err := c.ListByCommunicationsGateway(ctx, id)
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

	result = ListByCommunicationsGatewayCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
